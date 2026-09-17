package notifications

import (
	"database/sql"
	"errors"
	"social/internal/models"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestNotificationLifecycle(t *testing.T) {
	db := newNotificationTestDatabase(t)

	created, err := Create(db, 1, models.CreateNotificationRequest{
		Category: "requests",
		Type:     "follow_request",
		Message:  "Please review this request",
	})
	if err != nil {
		t.Fatal(err)
	}

	items, err := List(db, 1, "requests")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != created.ID || items[0].IsRead {
		t.Fatalf("unexpected notification list: %+v", items)
	}

	count, err := UnreadCount(db, 1)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("unread count = %d, expected 1", count)
	}

	if err := MarkRead(db, 1, created.ID); err != nil {
		t.Fatal(err)
	}
	count, err = UnreadCount(db, 1)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("unread count after mark read = %d, expected 0", count)
	}
}

func TestNotificationRejectsInvalidCategory(t *testing.T) {
	db := newNotificationTestDatabase(t)

	_, err := Create(db, 1, models.CreateNotificationRequest{
		Category: "unknown",
		Type:     "test",
		Message:  "invalid",
	})
	if !errors.Is(err, ErrInvalidCategory) {
		t.Fatalf("expected ErrInvalidCategory, got %v", err)
	}
}

func TestMessageNotificationIsSupported(t *testing.T) {
	db := newNotificationTestDatabase(t)

	created, err := Create(db, 1, models.CreateNotificationRequest{
		Category: "messages",
		Type:     "new_message",
		Message:  "A new message is waiting",
	})
	if err != nil {
		t.Fatal(err)
	}

	items, err := List(db, 1, "messages")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != created.ID || items[0].Category != "messages" {
		t.Fatalf("unexpected message notification list: %+v", items)
	}
}

func TestUnreadCountCanBeLimitedToMessages(t *testing.T) {
	db := newNotificationTestDatabase(t)

	_, err := db.Exec(`
		INSERT INTO notifications (user_id, category, type, message, is_read)
		VALUES
			(1, 'messages', 'new_message', 'Unread chat', 0),
			(1, 'groups', 'group_update', 'Unread group update', 0),
			(1, 'messages', 'new_message', 'Read chat', 1)
	`)
	if err != nil {
		t.Fatal(err)
	}

	count, err := UnreadCount(db, 1, "messages")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("message unread count = %d, expected 1", count)
	}
}

func TestNotificationsCanLoadOnePageAtATime(t *testing.T) {
	db := newNotificationTestDatabase(t)

	for _, message := range []string{"one", "two", "three"} {
		_, err := Create(db, 1, models.CreateNotificationRequest{
			Category: "groups",
			Type:     "group_update",
			Message:  message,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	page, err := List(db, 1, "groups", 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(page) != 2 {
		t.Fatalf("first notification page length = %d, expected 2", len(page))
	}

	olderPage, err := List(db, 1, "groups", 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(olderPage) != 1 {
		t.Fatalf("older notification page length = %d, expected 1", len(olderPage))
	}
}

func TestMessageNotificationsCanBeMarkedReadForOneChat(t *testing.T) {
	db := newNotificationTestDatabase(t)

	_, err := db.Exec(`
		INSERT INTO notifications (user_id, category, type, message, related_id, is_read)
		VALUES
			(1, 'messages', 'new_message', 'Chat one', 10, 0),
			(1, 'messages', 'new_message', 'Chat two', 11, 0),
			(1, 'groups', 'group_update', 'Keep this unread', 10, 0)
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := MarkMessageNotificationsRead(db, 1, 10); err != nil {
		t.Fatal(err)
	}

	var read, otherChat, otherCategory int
	err = db.QueryRow("SELECT is_read FROM notifications WHERE category = 'messages' AND related_id = 10").Scan(&read)
	if err != nil {
		t.Fatal(err)
	}
	err = db.QueryRow("SELECT is_read FROM notifications WHERE category = 'messages' AND related_id = 11").Scan(&otherChat)
	if err != nil {
		t.Fatal(err)
	}
	err = db.QueryRow("SELECT is_read FROM notifications WHERE category = 'groups' AND related_id = 10").Scan(&otherCategory)
	if err != nil {
		t.Fatal(err)
	}

	if read != 1 || otherChat != 0 || otherCategory != 0 {
		t.Fatalf("read states = %d, %d, %d; wanted 1, 0, 0", read, otherChat, otherCategory)
	}
}

func TestJoinRequestNotificationUsesSpecificRequest(t *testing.T) {
	db := newNotificationTestDatabase(t)

	rejected := "rejected"
	pending := "pending"

	_, err := db.Exec(`
		INSERT INTO user (id) VALUES (2);
		INSERT INTO group_join_requests (id, group_id, user_id, status)
		VALUES (10, 7, 2, 'rejected'), (11, 7, 2, 'pending');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES
			(20, 1, 2, 'groups', 'join_request', 'first request', 10),
			(21, 1, 2, 'groups', 'join_request', 'second request', 11);
	`)
	if err != nil {
		t.Fatal(err)
	}

	items, err := List(db, 1, "groups")
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Fatalf("notification count = %d, expected 2", len(items))
	}
	if items[0].ID != 21 || items[0].RequestStatus == nil || *items[0].RequestStatus != pending {
		t.Fatalf("new notification status = %+v, expected pending", items[0])
	}
	if items[1].ID != 20 || items[1].RequestStatus == nil || *items[1].RequestStatus != rejected {
		t.Fatalf("old notification status = %+v, expected rejected", items[1])
	}
}

func TestInvitationNotificationUsesSpecificInvitation(t *testing.T) {
	db := newNotificationTestDatabase(t)

	_, err := db.Exec(`
		INSERT INTO group_invitations (id, group_id, user_id, status)
		VALUES (10, 7, 1, 'declined'), (11, 7, 1, 'pending');
		INSERT INTO notifications (id, user_id, category, type, message, related_id)
		VALUES
			(20, 1, 'groups', 'invitation', 'old invitation', 10),
			(21, 1, 'groups', 'invitation', 'new invitation', 11);
	`)
	if err != nil {
		t.Fatal(err)
	}

	items, err := List(db, 1, "groups")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("notification count = %d, expected 2", len(items))
	}
	if items[0].InvitationStatus == nil || *items[0].InvitationStatus != "pending" {
		t.Fatalf("new invitation status = %+v", items[0].InvitationStatus)
	}
	if items[1].InvitationStatus == nil || *items[1].InvitationStatus != "declined" {
		t.Fatalf("old invitation status = %+v", items[1].InvitationStatus)
	}
}

func newNotificationTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE user (id INTEGER PRIMARY KEY);
		CREATE TABLE user_followers (target_id INTEGER, follower_id INTEGER, status INTEGER);
		CREATE TABLE event_rsvps (event_id INTEGER, user_id INTEGER, response TEXT);
		CREATE TABLE notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			actor_id INTEGER,
			category TEXT NOT NULL,
			type TEXT NOT NULL,
			message TEXT NOT NULL,
			related_id INTEGER,
			is_read INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id)
		);
		CREATE TABLE group_join_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE UNIQUE INDEX group_join_requests_one_pending
		ON group_join_requests (group_id, user_id)
		WHERE status = 'pending';
		CREATE TABLE group_invitations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		INSERT INTO user (id) VALUES (1);
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

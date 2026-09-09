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
		INSERT INTO user (id) VALUES (1);
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

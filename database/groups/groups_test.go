package groups

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestJoinGroupCreatesMembership(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE groups (id INTEGER PRIMARY KEY, title TEXT, description TEXT, creator_id INTEGER);
		CREATE TABLE chats (id INTEGER PRIMARY KEY AUTOINCREMENT, type TEXT, group_id INTEGER, num_of_members INTEGER DEFAULT 0);
		CREATE TABLE chat_users (user_id INTEGER, chat_id INTEGER, is_owner INTEGER DEFAULT 0, PRIMARY KEY (user_id, chat_id));
		CREATE TABLE group_members (group_id INTEGER, user_id INTEGER, joined_at DATETIME DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (group_id, user_id));
		CREATE TABLE group_invitations (id INTEGER PRIMARY KEY, group_id INTEGER, user_id INTEGER, status TEXT NOT NULL DEFAULT 'pending');
		INSERT INTO groups (id, title, description) VALUES (7, 'Orbit hikers', 'Walks');
		INSERT INTO chats (id, type, group_id) VALUES (11, 'group', 7);
		INSERT INTO group_invitations (id, group_id, user_id, status) VALUES (21, 7, 3, 'pending');
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := AcceptInvitation(db, 3, 21); err != nil {
		t.Fatal(err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM group_members").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("membership count = %d, want 1", count)
	}
}

func TestInvitationActionsUseExactInvitation(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE groups (id INTEGER PRIMARY KEY, creator_id INTEGER);
		CREATE TABLE chats (id INTEGER PRIMARY KEY, group_id INTEGER);
		CREATE TABLE chat_users (user_id INTEGER, chat_id INTEGER, PRIMARY KEY (user_id, chat_id));
		CREATE TABLE group_members (group_id INTEGER, user_id INTEGER, PRIMARY KEY (group_id, user_id));
		CREATE TABLE group_invitations (
			id INTEGER PRIMARY KEY,
			group_id INTEGER,
			user_id INTEGER,
			status TEXT NOT NULL
		);
		INSERT INTO groups (id, creator_id) VALUES (7, 1);
		INSERT INTO group_invitations (id, group_id, user_id, status) VALUES
			(10, 7, 3, 'declined'),
			(11, 7, 3, 'pending');
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := AcceptInvitation(db, 3, 10); !errors.Is(err, ErrInvitationNotFound) {
		t.Fatalf("historical invitation action error = %v, expected ErrInvitationNotFound", err)
	}
	if err := AcceptInvitation(db, 3, 11); err != nil {
		t.Fatal(err)
	}

	var oldStatus, newStatus string
	if err := db.QueryRow("SELECT status FROM group_invitations WHERE id = 10").Scan(&oldStatus); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT status FROM group_invitations WHERE id = 11").Scan(&newStatus); err != nil {
		t.Fatal(err)
	}
	if oldStatus != "declined" || newStatus != "accepted" {
		t.Fatalf("invitation statuses = %q, %q", oldStatus, newStatus)
	}

	if _, err := db.Exec(`
		INSERT INTO group_invitations (id, group_id, user_id, status) VALUES
			(12, 7, 4, 'accepted'),
			(13, 7, 4, 'pending')
	`); err != nil {
		t.Fatal(err)
	}
	if err := DeclineInvitation(db, 4, 12); !errors.Is(err, ErrInvitationNotFound) {
		t.Fatalf("historical decline error = %v, expected ErrInvitationNotFound", err)
	}
	if err := DeclineInvitation(db, 4, 13); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT status FROM group_invitations WHERE id = 13").Scan(&newStatus); err != nil {
		t.Fatal(err)
	}
	if newStatus != "declined" {
		t.Fatalf("declined invitation status = %q", newStatus)
	}
}

func TestAcceptJoinRequestUpdatesOnlyPendingRequest(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE groups (id INTEGER PRIMARY KEY, title TEXT, description TEXT, creator_id INTEGER);
		CREATE TABLE group_members (group_id INTEGER, user_id INTEGER, joined_at DATETIME DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (group_id, user_id));
		CREATE TABLE group_join_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE UNIQUE INDEX group_join_requests_one_pending
		ON group_join_requests (group_id, user_id)
		WHERE status = 'pending';
		INSERT INTO groups (id, title, description, creator_id) VALUES (7, 'Orbit hikers', 'Walks', 1);
		INSERT INTO group_join_requests (id, group_id, user_id, status)
		VALUES (10, 7, 2, 'rejected'), (11, 7, 2, 'pending');
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := AcceptJoinRequest(db, 1, 2, 7); err != nil {
		t.Fatal(err)
	}

	statuses := map[int]string{}
	rows, err := db.Query(`
		SELECT id, status
		FROM group_join_requests
		ORDER BY id
	`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var status string
		if err := rows.Scan(&id, &status); err != nil {
			t.Fatal(err)
		}
		statuses[id] = status
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	if statuses[10] != "rejected" {
		t.Fatalf("old request status = %q, expected rejected", statuses[10])
	}
	if statuses[11] != "accepted" {
		t.Fatalf("pending request status = %q, expected accepted", statuses[11])
	}

	var memberCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = 7 AND user_id = 2").Scan(&memberCount); err != nil {
		t.Fatal(err)
	}
	if memberCount != 1 {
		t.Fatalf("member count = %d, expected 1", memberCount)
	}
}

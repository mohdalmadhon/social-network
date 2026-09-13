package events

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSetRSVPUpdatesExistingResponse(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE user (id INTEGER PRIMARY KEY);
		CREATE TABLE group_members (group_id INTEGER, user_id INTEGER);
		INSERT INTO group_members VALUES (1,3);
		CREATE TABLE events (id INTEGER PRIMARY KEY, group_id INTEGER, creator_id INTEGER, title TEXT, content TEXT);
		CREATE TABLE event_rsvps (
			event_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			response TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (event_id, user_id)
		);
		INSERT INTO user (id) VALUES (3);
		INSERT INTO events (id, group_id, creator_id, title, content) VALUES (8, 1, 3, 'Meetup', 'Meetup');
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := SetRSVP(db, 3, 8, "going"); err != nil {
		t.Fatal(err)
	}
	if err := SetRSVP(db, 3, 8, "declined"); err != nil {
		t.Fatal(err)
	}

	var response string
	if err := db.QueryRow("SELECT response FROM event_rsvps WHERE event_id = 8 AND user_id = 3").Scan(&response); err != nil {
		t.Fatal(err)
	}
	if response != "declined" {
		t.Fatalf("response = %q, want declined", response)
	}
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM event_rsvps WHERE event_id = 8 AND user_id = 3").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("RSVP count = %d, want 1", count)
	}
}

func TestGroupRSVPRemovalAndEventDeletion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE user (id INTEGER PRIMARY KEY);
		CREATE TABLE groups (id INTEGER PRIMARY KEY);
		CREATE TABLE group_members (group_id INTEGER, user_id INTEGER, PRIMARY KEY (group_id, user_id));
		CREATE TABLE events (
			id INTEGER PRIMARY KEY,
			group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			creator_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
			title TEXT,
			content TEXT
		);
		CREATE TABLE event_rsvps (
			event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL,
			response TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (event_id, user_id)
		);
		CREATE TABLE notifications (id INTEGER PRIMARY KEY, category TEXT, related_id INTEGER);
		CREATE TRIGGER delete_event_notifications AFTER DELETE ON events BEGIN
			DELETE FROM notifications WHERE category = 'events' AND related_id = OLD.id;
		END;
		INSERT INTO user (id) VALUES (3), (4), (5);
		INSERT INTO groups (id) VALUES (1), (2);
		INSERT INTO group_members VALUES (1, 3), (1, 4), (2, 5);
		INSERT INTO events (id, group_id, creator_id, title, content) VALUES
			(8, 1, 3, 'Meetup', 'Meetup'),
			(9, 2, 5, 'Other', 'Other');
		INSERT INTO event_rsvps (event_id, user_id, response) VALUES (8, 3, 'going'), (8, 4, 'declined');
		INSERT INTO notifications (id, category, related_id) VALUES (20, 'events', 8), (21, 'events', 9);
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := SetGroupRSVP(db, 4, 1, 8, "going"); err != nil {
		t.Fatal(err)
	}
	var response string
	if err := db.QueryRow(`SELECT response FROM event_rsvps WHERE event_id = 8 AND user_id = 4`).Scan(&response); err != nil {
		t.Fatal(err)
	}
	if response != "going" {
		t.Fatalf("switched response = %q, want going", response)
	}
	if err := SetGroupRSVP(db, 5, 2, 8, "going"); !errors.Is(err, ErrEventNotFound) {
		t.Fatalf("mismatched group error = %v, want ErrEventNotFound", err)
	}

	if err := RemoveRSVP(db, 4, 1, 8); err != nil {
		t.Fatal(err)
	}
	if err := RemoveRSVP(db, 4, 1, 8); err != nil {
		t.Fatalf("removing an absent RSVP should be idempotent: %v", err)
	}
	assertEventRowCount(t, db, `SELECT COUNT(*) FROM event_rsvps WHERE event_id = 8 AND user_id = 4`, 0)
	assertEventRowCount(t, db, `SELECT COUNT(*) FROM event_rsvps WHERE event_id = 8 AND user_id = 3`, 1)

	if err := Delete(db, 4, 1, 8); !errors.Is(err, ErrNotCreator) {
		t.Fatalf("member delete error = %v, want ErrNotCreator", err)
	}
	if err := Delete(db, 3, 2, 8); !errors.Is(err, ErrEventNotFound) {
		t.Fatalf("mismatched delete error = %v, want ErrEventNotFound", err)
	}
	if err := Delete(db, 3, 1, 8); err != nil {
		t.Fatal(err)
	}
	assertEventRowCount(t, db, `SELECT COUNT(*) FROM events WHERE id = 8`, 0)
	assertEventRowCount(t, db, `SELECT COUNT(*) FROM event_rsvps WHERE event_id = 8`, 0)
	assertEventRowCount(t, db, `SELECT COUNT(*) FROM notifications WHERE id = 20`, 0)
	assertEventRowCount(t, db, `SELECT COUNT(*) FROM notifications WHERE id = 21`, 1)
}

func assertEventRowCount(t *testing.T, db *sql.DB, query string, expected int) {
	t.Helper()
	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != expected {
		t.Fatalf("row count = %d, expected %d for %q", count, expected, query)
	}
}

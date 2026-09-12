package events

import (
	"database/sql"
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
}

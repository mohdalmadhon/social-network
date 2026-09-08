package groups

import (
	"database/sql"
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
		CREATE TABLE groups (id INTEGER PRIMARY KEY, title TEXT, description TEXT);
		CREATE TABLE chats (id INTEGER PRIMARY KEY AUTOINCREMENT, type TEXT, group_id INTEGER, num_of_members INTEGER DEFAULT 0);
		CREATE TABLE chat_users (user_id INTEGER, chat_id INTEGER, is_owner INTEGER DEFAULT 0, PRIMARY KEY (user_id, chat_id));
		CREATE TABLE group_members (group_id INTEGER, user_id INTEGER, joined_at DATETIME DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (group_id, user_id));
		INSERT INTO groups (id, title, description) VALUES (7, 'Orbit hikers', 'Walks');
		INSERT INTO chats (id, type, group_id) VALUES (11, 'group', 7);
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := AcceptInvitation(db, 3, 7); err != nil {
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

package chats

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestPrivateChatsRequireFollowRelationshipAndReusePair(t *testing.T) {
	db := newChatTestDB(t)

	first, err := OpenPrivateChat(db, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	second, err := OpenPrivateChat(db, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID {
		t.Fatalf("reverse open created chat %d; expected %d", second.ID, first.ID)
	}
	assertChatCount(t, db, "SELECT COUNT(*) FROM chats WHERE type = 'private'", 1)
	assertChatCount(t, db, "SELECT COUNT(*) FROM chat_users WHERE chat_id = ?", 2, first.ID)

	if _, err = OpenPrivateChat(db, 1, 3); !errors.Is(err, ErrNoFollowRelation) {
		t.Fatalf("unrelated users error = %v, expected %v", err, ErrNoFollowRelation)
	}
	if _, err = OpenPrivateChat(db, 1, 1); !errors.Is(err, ErrSelfChat) {
		t.Fatalf("self chat error = %v, expected %v", err, ErrSelfChat)
	}
	if _, err = OpenPrivateChat(db, 1, 99); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("missing user error = %v, expected %v", err, ErrUserNotFound)
	}

	candidates, err := ListPrivateCandidates(db, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 || candidates[0].ID != 1 || candidates[0].ChatID == nil || *candidates[0].ChatID != first.ID {
		t.Fatalf("candidate list = %+v", candidates)
	}
}

func TestPrivateMessagesEnforceParticipantsRelationshipAndChronology(t *testing.T) {
	db := newChatTestDB(t)
	conversation, err := OpenPrivateChat(db, 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	first, err := SendPrivateMessage(db, conversation.ID, 1, "Hello 👋😂❤️")
	if err != nil {
		t.Fatal(err)
	}
	second, err := SendPrivateMessage(db, conversation.ID, 2, "  Welcome back  ")
	if err != nil {
		t.Fatal(err)
	}
	if first.Content != "Hello 👋😂❤️" || second.Content != "Welcome back" {
		t.Fatalf("message content changed: %q / %q", first.Content, second.Content)
	}

	messages, err := ListPrivateMessages(db, conversation.ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 || messages[0].ID != first.ID || messages[1].ID != second.ID {
		t.Fatalf("messages are not chronological: %+v", messages)
	}
	if !messages[0].IsOwn || messages[1].IsOwn {
		t.Fatalf("message ownership is incorrect: %+v", messages)
	}
	latestPage, err := ListPrivateMessages(db, conversation.ID, 1, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(latestPage) != 1 || latestPage[0].ID != second.ID {
		t.Fatalf("latest message page = %+v, expected the newest message", latestPage)
	}
	olderPage, err := ListPrivateMessages(db, conversation.ID, 1, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(olderPage) != 1 || olderPage[0].ID != first.ID {
		t.Fatalf("older message page = %+v, expected the first message", olderPage)
	}
	if _, err = ListPrivateMessages(db, conversation.ID, 3); !errors.Is(err, ErrForbidden) {
		t.Fatalf("non-participant error = %v, expected %v", err, ErrForbidden)
	}
	if _, err = SendPrivateMessage(db, conversation.ID, 1, "   "); !errors.Is(err, ErrInvalidMessage) {
		t.Fatalf("empty message error = %v, expected %v", err, ErrInvalidMessage)
	}

	if _, err = db.Exec(`DELETE FROM user_followers WHERE follower_id = 1 AND target_id = 2`); err != nil {
		t.Fatal(err)
	}
	if _, err = SendPrivateMessage(db, conversation.ID, 2, "blocked"); !errors.Is(err, ErrNoFollowRelation) {
		t.Fatalf("removed relationship send error = %v, expected %v", err, ErrNoFollowRelation)
	}
	if _, err = ListPrivateMessages(db, conversation.ID, 1); !errors.Is(err, ErrNoFollowRelation) {
		t.Fatalf("removed relationship read error = %v, expected %v", err, ErrNoFollowRelation)
	}
}

func TestGroupChatRequiresMembershipAndUsesOneRoom(t *testing.T) {
	db := newChatTestDB(t)

	first, err := SendGroupMessage(db, 7, 1, "Welcome to the group")
	if err != nil {
		t.Fatal(err)
	}
	second, err := SendGroupMessage(db, 7, 2, "Hello everyone 🔥")
	if err != nil {
		t.Fatal(err)
	}
	messages, err := ListGroupMessages(db, 7, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 || messages[0].ID != first.ID || messages[1].ID != second.ID {
		t.Fatalf("group messages = %+v", messages)
	}
	assertChatCount(t, db, "SELECT COUNT(*) FROM chats WHERE type = 'group' AND group_id = 7", 1)

	if _, err = ListGroupMessages(db, 7, 3); !errors.Is(err, ErrForbidden) {
		t.Fatalf("non-member read error = %v, expected %v", err, ErrForbidden)
	}
	if _, err = SendGroupMessage(db, 7, 3, "not allowed"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("non-member send error = %v, expected %v", err, ErrForbidden)
	}
	if _, err = ListGroupMessages(db, 99, 1); !errors.Is(err, ErrGroupNotFound) {
		t.Fatalf("missing group error = %v, expected %v", err, ErrGroupNotFound)
	}
}

func newChatTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE user (
			id INTEGER PRIMARY KEY,
			username TEXT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL
		);
		CREATE TABLE profile (user_id INTEGER PRIMARY KEY, avatar_path TEXT);
		CREATE TABLE user_followers (
			follower_id INTEGER NOT NULL,
			target_id INTEGER NOT NULL,
			status INTEGER NOT NULL,
			PRIMARY KEY (follower_id, target_id)
		);
		CREATE TABLE groups (id INTEGER PRIMARY KEY, title TEXT NOT NULL);
		CREATE TABLE group_members (
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (group_id, user_id)
		);
		CREATE TABLE chats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			num_of_members INTEGER NOT NULL DEFAULT 0,
			private_user_low_id INTEGER REFERENCES user(id) ON DELETE CASCADE,
			private_user_high_id INTEGER REFERENCES user(id) ON DELETE CASCADE
		);
		CREATE UNIQUE INDEX chats_one_room_per_group ON chats(group_id) WHERE type = 'group' AND group_id IS NOT NULL;
		CREATE UNIQUE INDEX chats_one_private_room_per_pair ON chats(private_user_low_id, private_user_high_id)
			WHERE type = 'private' AND private_user_low_id IS NOT NULL AND private_user_high_id IS NOT NULL;
		CREATE TABLE chat_users (
			user_id INTEGER NOT NULL,
			chat_id INTEGER NOT NULL,
			is_owner INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (user_id, chat_id)
		);
		CREATE TABLE messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sender_id INTEGER NOT NULL,
			chat_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		INSERT INTO user (id, username, first_name, last_name) VALUES
			(1, 'alice', 'Alice', 'Orbit'),
			(2, 'bob', 'Bob', 'Orbit'),
			(3, 'casey', 'Casey', 'Orbit');
		INSERT INTO profile (user_id, avatar_path) VALUES
			(1, 'avatars/alice.png'), (2, 'avatars/bob.png'), (3, '');
		INSERT INTO user_followers (follower_id, target_id, status) VALUES
			(1, 2, 1), (1, 3, 0);
		INSERT INTO groups (id, title) VALUES (7, 'Design Guild');
		INSERT INTO group_members (group_id, user_id) VALUES (7, 1), (7, 2);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func assertChatCount(t *testing.T, db *sql.DB, query string, expected int, args ...any) {
	t.Helper()
	var count int
	if err := db.QueryRow(query, args...).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != expected {
		t.Fatalf("count = %d, expected %d for %q", count, expected, query)
	}
}

package migrations

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestPostPrivacyMigrationUpDownUp(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "002_init_posts.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "004_create_triggers.up.sql")
	insertLegacyPostData(t, db)
	runMigrationFile(t, db, "005_post_privacy.up.sql")

	assertColumn(t, db, "posts", "privacy", true)
	assertColumnNotNull(t, db, "posts", "group_id", false)
	assertTable(t, db, "post_viewers", true)
	assertLegacyPostData(t, db)

	runMigrationFile(t, db, "005_post_privacy.down.sql")
	assertColumn(t, db, "posts", "privacy", false)
	assertTable(t, db, "post_viewers", false)

	runMigrationFile(t, db, "005_post_privacy.up.sql")
	assertColumn(t, db, "posts", "privacy", true)
	assertTable(t, db, "post_viewers", true)
	assertLegacyPostData(t, db)
}

func TestNotificationsMigrationUpDownUp(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "006_notifications.up.sql")

	_, err = db.Exec(`
		INSERT INTO user (id, email, username, first_name, last_name, dob, password)
		VALUES (1, 'notification@orbit.test', 'notification', 'Notification', 'User', '2000-01-01', 'password');
		INSERT INTO notifications (user_id, category, type, message, related_id)
		VALUES (1, 'requests', 'follow_request', 'Someone wants to follow you', 42);
	`)
	if err != nil {
		t.Fatal(err)
	}

	var message string
	if err := db.QueryRow("SELECT message FROM notifications WHERE user_id = 1").Scan(&message); err != nil {
		t.Fatal(err)
	}
	if message != "Someone wants to follow you" {
		t.Fatalf("notification message changed: %q", message)
	}

	runMigrationFile(t, db, "006_notifications.down.sql")
	assertTable(t, db, "notifications", false)
	runMigrationFile(t, db, "006_notifications.up.sql")
	assertTable(t, db, "notifications", true)
}

func TestCommentMediaMigrationUpDownUp(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "002_init_posts.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "004_create_triggers.up.sql")
	runMigrationFile(t, db, "005_post_privacy.up.sql")
	runMigrationFile(t, db, "007_comment_media.up.sql")

	assertColumn(t, db, "comments", "image_path", true)

	runMigrationFile(t, db, "007_comment_media.down.sql")
	assertColumn(t, db, "comments", "image_path", false)

	runMigrationFile(t, db, "007_comment_media.up.sql")
	assertColumn(t, db, "comments", "image_path", true)
}

func TestEventRsvpMigrationUpDownUp(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "008_event_rsvps.up.sql")

	assertTable(t, db, "event_rsvps", true)
	runMigrationFile(t, db, "008_event_rsvps.down.sql")
	assertTable(t, db, "event_rsvps", false)
	runMigrationFile(t, db, "008_event_rsvps.up.sql")
	assertTable(t, db, "event_rsvps", true)
}

func TestGroupMigrationsUpDownUp(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "002_init_posts.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "004_create_triggers.up.sql")
	runMigrationFile(t, db, "005_post_privacy.up.sql")
	runMigrationFile(t, db, "006_notifications.up.sql")
	runMigrationFile(t, db, "007_comment_media.up.sql")
	runMigrationFile(t, db, "008_event_rsvps.up.sql")
	runMigrationFile(t, db, "009_add_creator_id_to_groups.up.sql")
	runMigrationFile(t, db, "010_create_group_members.up.sql")

	assertColumn(t, db, "groups", "creator_id", true)
	assertTable(t, db, "group_members", true)

	runMigrationFile(t, db, "010_create_group_members.down.sql")
	assertTable(t, db, "group_members", false)
	runMigrationFile(t, db, "010_create_group_members.up.sql")
	assertTable(t, db, "group_members", true)

	runMigrationFile(t, db, "009_add_creator_id_to_groups.down.sql")
	assertColumn(t, db, "groups", "creator_id", false)
	runMigrationFile(t, db, "009_add_creator_id_to_groups.up.sql")
	assertColumn(t, db, "groups", "creator_id", true)
}

func insertLegacyPostData(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		INSERT INTO user (id, email, username, first_name, last_name, dob, password)
		VALUES (1, 'legacy@orbit.test', 'legacy', 'Legacy', 'User', '2000-01-01', 'password');

		INSERT INTO groups (id, title, description)
		VALUES (1, 'Legacy Group', 'Migration test group');

		INSERT INTO posts (id, type, title, content, image_path, user_id, group_id)
		VALUES (1, 'group', 'Legacy title', 'Legacy post', '', 1, 1);

		INSERT INTO comments (id, user_id, post_id, content)
		VALUES (1, 1, 1, 'Legacy comment');

		INSERT INTO post_reactions (user_id, post_id, value)
		VALUES (1, 1, 1);
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func assertLegacyPostData(t *testing.T, db *sql.DB) {
	t.Helper()

	var content, privacy string
	var groupID, likes, comments int
	err := db.QueryRow(`
		SELECT content, privacy, group_id, like_count, comment_count
		FROM posts
		WHERE id = 1
	`).Scan(&content, &privacy, &groupID, &likes, &comments)
	if err != nil {
		t.Fatal(err)
	}

	if content != "Legacy post" || privacy != "public" || groupID != 1 || likes != 1 || comments != 1 {
		t.Fatalf("legacy post data changed: content=%q privacy=%q group=%d likes=%d comments=%d", content, privacy, groupID, likes, comments)
	}
}

func runMigrationFile(t *testing.T, db *sql.DB, filename string) {
	t.Helper()

	query, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(query)); err != nil {
		t.Fatalf("%s failed: %v", filename, err)
	}
}

func assertColumn(t *testing.T, db *sql.DB, table string, column string, expected bool) {
	t.Helper()

	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id int
		var name, dataType string
		var notNull, primaryKey int
		var defaultValue any
		if err = rows.Scan(&id, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatal(err)
		}
		if name == column {
			found = true
		}
	}

	if found != expected {
		t.Fatalf("column %s.%s found=%v, expected %v", table, column, found, expected)
	}
}

func assertColumnNotNull(t *testing.T, db *sql.DB, table string, column string, expected bool) {
	t.Helper()

	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, dataType string
		var notNull, primaryKey int
		var defaultValue any
		if err = rows.Scan(&id, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatal(err)
		}
		if name == column {
			if (notNull == 1) != expected {
				t.Fatalf("column %s.%s not-null=%v, expected %v", table, column, notNull == 1, expected)
			}
			return
		}
	}

	t.Fatalf("column %s.%s was not found", table, column)
}

func assertTable(t *testing.T, db *sql.DB, table string, expected bool) {
	t.Helper()

	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM sqlite_master
		WHERE type = 'table' AND name = ?
	`, table).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if (count == 1) != expected {
		t.Fatalf("table %s found=%v, expected %v", table, count == 1, expected)
	}
}

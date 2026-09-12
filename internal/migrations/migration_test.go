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

func TestGroupJoinRequestHistoryMigrationPreservesData(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "006_notifications.up.sql")
	runMigrationFile(t, db, "009_add_creator_id_to_groups.up.sql")
	runMigrationFile(t, db, "011_create_group_join_requests.up.sql")

	_, err = db.Exec(`
		INSERT INTO user (id, email, username, first_name, last_name, dob, password)
		VALUES
			(1, 'owner@orbit.test', 'owner', 'Group', 'Owner', '2000-01-01', 'password'),
			(2, 'requester@orbit.test', 'requester', 'Group', 'Requester', '2000-01-01', 'password');
		INSERT INTO groups (id, title, description, creator_id)
		VALUES (7, 'Orbit hikers', 'Walks', 1);
		INSERT INTO group_join_requests (group_id, user_id, status)
		VALUES (7, 2, 'rejected');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES (20, 1, 2, 'groups', 'join_request', 'requested to join your group', 7);
	`)
	if err != nil {
		t.Fatal(err)
	}

	runMigrationFile(t, db, "014_group_join_request_history.up.sql")

	assertColumn(t, db, "group_join_requests", "id", true)

	var requestID int64
	var status string
	err = db.QueryRow(`
		SELECT id, status
		FROM group_join_requests
		WHERE group_id = 7 AND user_id = 2
	`).Scan(&requestID, &status)
	if err != nil {
		t.Fatal(err)
	}
	if status != "rejected" {
		t.Fatalf("status = %q, expected rejected", status)
	}

	var relatedID int64
	err = db.QueryRow("SELECT related_id FROM notifications WHERE id = 20").Scan(&relatedID)
	if err != nil {
		t.Fatal(err)
	}
	if relatedID != requestID {
		t.Fatalf("notification related_id = %d, expected request id %d", relatedID, requestID)
	}

	_, err = db.Exec(`
		INSERT INTO group_join_requests (group_id, user_id, status)
		VALUES (7, 2, 'pending')
	`)
	if err != nil {
		t.Fatalf("new pending request after rejection failed: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO group_join_requests (group_id, user_id, status)
		VALUES (7, 2, 'pending')
	`)
	if err == nil {
		t.Fatal("duplicate pending request succeeded")
	}
}

func TestGroupDeletionCleansUpJoinRequestNotifications(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatal(err)
	}

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "006_notifications.up.sql")
	runMigrationFile(t, db, "009_add_creator_id_to_groups.up.sql")
	runMigrationFile(t, db, "011_create_group_join_requests.up.sql")
	runMigrationFile(t, db, "013_add_notification_cleanup_triggers.up.sql")
	runMigrationFile(t, db, "014_group_join_request_history.up.sql")
	runMigrationFile(t, db, "015_fix_group_notification_cleanup.up.sql")

	_, err = db.Exec(`
		INSERT INTO user (id, email, username, first_name, last_name, dob, password)
		VALUES
			(1, 'owner@orbit.test', 'owner', 'Group', 'Owner', '2000-01-01', 'password'),
			(2, 'requester@orbit.test', 'requester', 'Group', 'Requester', '2000-01-01', 'password'),
			(3, 'other@orbit.test', 'other', 'Other', 'Owner', '2000-01-01', 'password');
		INSERT INTO groups (id, title, description, creator_id)
		VALUES
			(7, 'Orbit hikers', 'Walks', 1),
			(8, 'Orbit readers', 'Books', 3);
		INSERT INTO group_join_requests (id, group_id, user_id, status)
		VALUES
			(10, 7, 2, 'rejected'),
			(11, 7, 2, 'pending'),
			(7, 8, 2, 'rejected');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES
			(20, 1, 2, 'groups', 'join_request', 'old request', 10),
			(21, 1, 2, 'groups', 'join_request', 'pending request', 11),
			(22, 2, 1, 'groups', 'group_invitation', 'group invitation', 7),
			(23, 3, 2, 'groups', 'join_request', 'unrelated request', 7),
			(24, 2, 1, 'requests', 'follow_request', 'unrelated notification', 7);
	`)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec("DELETE FROM groups WHERE id = 7"); err != nil {
		t.Fatal(err)
	}

	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM group_join_requests WHERE group_id = 7", 0)
	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id IN (20, 21, 22)", 0)
	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id IN (23, 24)", 2)

	runMigrationFile(t, db, "015_fix_group_notification_cleanup.down.sql")
	runMigrationFile(t, db, "015_fix_group_notification_cleanup.up.sql")
}

func TestGroupInvitationHistoryMigrationPreservesData(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "006_notifications.up.sql")
	runMigrationFile(t, db, "009_add_creator_id_to_groups.up.sql")
	runMigrationFile(t, db, "011_create_group_join_requests.up.sql")
	runMigrationFile(t, db, "012_create_group_invitations.up.sql")
	runMigrationFile(t, db, "013_add_notification_cleanup_triggers.up.sql")
	runMigrationFile(t, db, "014_group_join_request_history.up.sql")
	runMigrationFile(t, db, "015_fix_group_notification_cleanup.up.sql")

	_, err = db.Exec(`
		INSERT INTO user (id, email, username, first_name, last_name, dob, password)
		VALUES
			(1, 'owner@orbit.test', 'owner', 'Group', 'Owner', '2000-01-01', 'password'),
			(2, 'invitee@orbit.test', 'invitee', 'Group', 'Invitee', '2000-01-01', 'password'),
			(3, 'member@orbit.test', 'member', 'Group', 'Member', '2000-01-01', 'password');
		INSERT INTO groups (id, title, description, creator_id)
		VALUES (7, 'Orbit hikers', 'Walks', 1);
		INSERT INTO group_invitations (group_id, user_id, inviter_id, status)
		VALUES (7, 2, 1, 'declined');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES (20, 2, 1, 'groups', 'invitation', 'invited you to join a group', 7);
	`)
	if err != nil {
		t.Fatal(err)
	}

	runMigrationFile(t, db, "016_group_invitation_history.up.sql")
	assertColumn(t, db, "group_invitations", "id", true)

	var invitationID int64
	var inviterID int
	var status string
	if err = db.QueryRow(`
		SELECT id, inviter_id, status
		FROM group_invitations
		WHERE group_id = 7 AND user_id = 2
	`).Scan(&invitationID, &inviterID, &status); err != nil {
		t.Fatal(err)
	}
	if inviterID != 1 || status != "declined" {
		t.Fatalf("migrated invitation inviter = %d, status = %q", inviterID, status)
	}

	var relatedID int64
	if err = db.QueryRow("SELECT related_id FROM notifications WHERE id = 20").Scan(&relatedID); err != nil {
		t.Fatal(err)
	}
	if relatedID != invitationID {
		t.Fatalf("notification related_id = %d, expected invitation id %d", relatedID, invitationID)
	}

	runMigrationFile(t, db, "016_group_invitation_history.down.sql")
	assertColumn(t, db, "group_invitations", "id", false)
	if err = db.QueryRow("SELECT related_id FROM notifications WHERE id = 20").Scan(&relatedID); err != nil {
		t.Fatal(err)
	}
	if relatedID != 7 {
		t.Fatalf("down migration related_id = %d, expected group id 7", relatedID)
	}

	runMigrationFile(t, db, "016_group_invitation_history.up.sql")
	if err = db.QueryRow(`
		SELECT id
		FROM group_invitations
		WHERE group_id = 7 AND user_id = 2 AND status = 'declined'
	`).Scan(&invitationID); err != nil {
		t.Fatal(err)
	}
	if err = db.QueryRow("SELECT related_id FROM notifications WHERE id = 20").Scan(&relatedID); err != nil {
		t.Fatal(err)
	}
	if relatedID != invitationID {
		t.Fatalf("reapplied migration related_id = %d, expected invitation id %d", relatedID, invitationID)
	}

	_, err = db.Exec(`
		INSERT INTO group_invitations (group_id, user_id, inviter_id, status)
		VALUES
			(7, 2, 3, 'declined'),
			(7, 2, 3, 'pending');
	`)
	if err != nil {
		t.Fatalf("historical invitations failed: %v", err)
	}
	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE group_id = 7 AND user_id = 2", 3)

	if _, err = db.Exec(`
		INSERT INTO group_invitations (group_id, user_id, inviter_id, status)
		VALUES (7, 2, 1, 'pending')
	`); err == nil {
		t.Fatal("duplicate pending invitation succeeded")
	}
}

func TestGroupInvitationNotificationCleanup(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatal(err)
	}

	runMigrationFile(t, db, "001_init_users.up.sql")
	runMigrationFile(t, db, "003_init_chats.up.sql")
	runMigrationFile(t, db, "006_notifications.up.sql")
	runMigrationFile(t, db, "009_add_creator_id_to_groups.up.sql")
	runMigrationFile(t, db, "011_create_group_join_requests.up.sql")
	runMigrationFile(t, db, "012_create_group_invitations.up.sql")
	runMigrationFile(t, db, "013_add_notification_cleanup_triggers.up.sql")
	runMigrationFile(t, db, "014_group_join_request_history.up.sql")
	runMigrationFile(t, db, "015_fix_group_notification_cleanup.up.sql")
	runMigrationFile(t, db, "016_group_invitation_history.up.sql")

	_, err = db.Exec(`
		INSERT INTO user (id, email, username, first_name, last_name, dob, password)
		VALUES
			(1, 'owner@orbit.test', 'owner', 'Group', 'Owner', '2000-01-01', 'password'),
			(2, 'invitee@orbit.test', 'invitee', 'Group', 'Invitee', '2000-01-01', 'password'),
			(3, 'other@orbit.test', 'other', 'Other', 'Owner', '2000-01-01', 'password');
		INSERT INTO groups (id, title, description, creator_id)
		VALUES
			(7, 'Orbit hikers', 'Walks', 1),
			(8, 'Orbit readers', 'Books', 3);
		INSERT INTO group_join_requests (id, group_id, user_id, status)
		VALUES (30, 7, 2, 'pending'), (7, 8, 2, 'rejected');
		INSERT INTO group_invitations (id, group_id, user_id, inviter_id, status)
		VALUES
			(10, 7, 2, 1, 'declined'),
			(11, 7, 2, 1, 'pending'),
			(7, 8, 2, 3, 'declined');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES
			(20, 2, 1, 'groups', 'invitation', 'old invitation', 10),
			(21, 2, 1, 'groups', 'invitation', 'pending invitation', 11),
			(22, 1, 2, 'groups', 'join_request', 'join request', 30),
			(23, 2, 1, 'groups', 'group_updated', 'group update', 7),
			(24, 2, 3, 'groups', 'invitation', 'unrelated invitation', 7),
			(25, 3, 2, 'groups', 'join_request', 'unrelated join request', 7),
			(26, 2, 1, 'requests', 'follow_request', 'unrelated notification', 7);
	`)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = db.Exec("DELETE FROM group_invitations WHERE id = 10"); err != nil {
		t.Fatal(err)
	}
	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 20", 0)
	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 21", 1)

	if _, err = db.Exec("DELETE FROM groups WHERE id = 7"); err != nil {
		t.Fatal(err)
	}
	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id IN (21, 22, 23)", 0)
	assertMigrationRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id IN (24, 25, 26)", 3)
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

func assertMigrationRowCount(t *testing.T, db *sql.DB, query string, expected int) {
	t.Helper()

	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != expected {
		t.Fatalf("row count = %d, expected %d for %q", count, expected, query)
	}
}

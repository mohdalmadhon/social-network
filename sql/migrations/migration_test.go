package migrations

import (
	"database/sql"
	"fmt"
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

	for version := 1; version <= 4; version++ {
		runMigrationFile(t, db, fmt.Sprintf("%06d_%s.up.sql", version, migrationName(version)))
	}
	insertLegacyPostData(t, db)
	runMigrationFile(t, db, "000005_post_privacy.up.sql")

	assertColumn(t, db, "posts", "privacy", true)
	assertColumnNotNull(t, db, "posts", "group_id", false)
	assertTable(t, db, "post_viewers", true)
	assertLegacyPostData(t, db)

	runMigrationFile(t, db, "000005_post_privacy.down.sql")
	assertColumn(t, db, "posts", "privacy", false)
	assertTable(t, db, "post_viewers", false)

	runMigrationFile(t, db, "000005_post_privacy.up.sql")
	assertColumn(t, db, "posts", "privacy", true)
	assertTable(t, db, "post_viewers", true)
	assertLegacyPostData(t, db)
}

func insertLegacyPostData(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		INSERT INTO users (id, email, username, first_name, last_name, dob, password)
		VALUES (1, 'legacy@orbit.test', 'legacy', 'Legacy', 'User', '2000-01-01', 'password');

		INSERT INTO profile (user_id, about)
		VALUES (1, 'legacy profile');

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

func migrationName(version int) string {
	names := map[int]string{
		1: "init",
		2: "triggers",
		3: "indexes",
		4: "update_users",
		5: "post_privacy",
	}
	return names[version]
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

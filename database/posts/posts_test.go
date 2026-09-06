package posts

import (
	"database/sql"
	"errors"
	"social/internal/models"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestPostPrivacyFiltersTheFeed(t *testing.T) {
	db := newPostTestDatabase(t)

	insertPostTestUser(t, db, 1, "author")
	insertPostTestUser(t, db, 2, "follower")
	insertPostTestUser(t, db, 3, "outsider")
	insertPostTestUser(t, db, 4, "selected")

	_, err := db.Exec(`
		INSERT INTO user_followers (follower_id, target_id, status)
		VALUES (2, 1, 1), (4, 1, 1)
	`)
	if err != nil {
		t.Fatal(err)
	}

	createTestPost(t, db, models.CreatePostRequest{
		Content: "public post",
		Privacy: models.PostPrivacyPublic,
	})
	createTestPost(t, db, models.CreatePostRequest{
		Content: "followers post",
		Privacy: models.PostPrivacyFollowers,
	})
	createTestPost(t, db, models.CreatePostRequest{
		Content:             "selected post",
		Privacy:             models.PostPrivacySelected,
		SelectedFollowerIDs: []int{4, 4},
	})

	assertFeedContents(t, db, 1, "public post", "followers post", "selected post")
	assertFeedContents(t, db, 2, "public post", "followers post")
	assertFeedContents(t, db, 3, "public post")
	assertFeedContents(t, db, 4, "public post", "followers post", "selected post")
}

func TestSelectedPostRejectsSomeoneWhoIsNotAFollower(t *testing.T) {
	db := newPostTestDatabase(t)
	insertPostTestUser(t, db, 1, "author")
	insertPostTestUser(t, db, 2, "outsider")

	_, err := CreatePost(db, 1, models.CreatePostRequest{
		Content:             "private post",
		Privacy:             models.PostPrivacySelected,
		SelectedFollowerIDs: []int{2},
	})
	if !errors.Is(err, ErrInvalidPostViewer) {
		t.Fatalf("expected ErrInvalidPostViewer, got %v", err)
	}
}

func createTestPost(t *testing.T, db *sql.DB, request models.CreatePostRequest) {
	t.Helper()
	if _, err := CreatePost(db, 1, request); err != nil {
		t.Fatal(err)
	}
}

func assertFeedContents(t *testing.T, db *sql.DB, viewerID int, expected ...string) {
	t.Helper()

	posts, err := ListFeedPosts(db, viewerID)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != len(expected) {
		t.Fatalf("viewer %d received %d posts, expected %d", viewerID, len(posts), len(expected))
	}

	found := make(map[string]bool)
	for _, post := range posts {
		found[post.Content] = true
	}
	for _, content := range expected {
		if !found[content] {
			t.Fatalf("viewer %d did not receive %q", viewerID, content)
		}
	}
}

func newPostTestDatabase(t *testing.T) *sql.DB {
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
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			dob DATE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE groups (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL
		);

		CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL DEFAULT 'post',
			title TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL,
			image_path TEXT NOT NULL DEFAULT '',
			user_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			group_id INTEGER,
			privacy TEXT NOT NULL DEFAULT 'public',
			like_count INTEGER NOT NULL DEFAULT 0,
			dislike_count INTEGER NOT NULL DEFAULT 0,
			comment_count INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES user(id),
			FOREIGN KEY (group_id) REFERENCES groups(id)
		);

		CREATE TABLE user_followers (
			follower_id INTEGER NOT NULL,
			target_id INTEGER NOT NULL,
			status INTEGER NOT NULL DEFAULT 1,
			PRIMARY KEY (follower_id, target_id),
			FOREIGN KEY (follower_id) REFERENCES user(id),
			FOREIGN KEY (target_id) REFERENCES user(id)
		);

		CREATE TABLE profile (
			user_id INTEGER PRIMARY KEY,
			avatar_path TEXT,
			FOREIGN KEY (user_id) REFERENCES user(id)
		);

		CREATE TABLE post_viewers (
			post_id INTEGER NOT NULL,
			viewer_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, viewer_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (viewer_id) REFERENCES user(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func insertPostTestUser(t *testing.T, db *sql.DB, id int, username string) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO user (id, email, username, first_name, last_name, dob, password)
		VALUES (?, ?, ?, 'Test', 'User', '2000-01-01', 'password')
	`, id, username+"@orbit.test", username)
	if err != nil {
		t.Fatal(err)
	}
}

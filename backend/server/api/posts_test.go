package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestPostsRequireAuthentication(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	response := httptest.NewRecorder()

	App{}.Posts(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestCreateAndListPublicPost(t *testing.T) {
	db := newAPIPostTestDatabase(t)
	t.Setenv("ORBIT_TOKEN_SECRET", "test-only-secret")
	token, err := GenerateToken(1, "noa")
	if err != nil {
		t.Fatal(err)
	}

	app := App{DB: db}
	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/posts",
		strings.NewReader(`{"content":"Hello, Orbit!","privacy":"public"}`),
	)
	createRequest.AddCookie(&http.Cookie{Name: "token", Value: token})
	createResponse := httptest.NewRecorder()
	app.Posts(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", createResponse.Code, createResponse.Body.String())
	}
	if !strings.Contains(createResponse.Body.String(), `"content":"Hello, Orbit!"`) {
		t.Fatalf("created post missing from response: %s", createResponse.Body.String())
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	listRequest.AddCookie(&http.Cookie{Name: "token", Value: token})
	listResponse := httptest.NewRecorder()
	app.Posts(listResponse, listRequest)

	if listResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listResponse.Code, listResponse.Body.String())
	}
	if !strings.Contains(listResponse.Body.String(), `"content":"Hello, Orbit!"`) {
		t.Fatalf("feed response missing post: %s", listResponse.Body.String())
	}
}

func TestCreatePostRejectsUnknownPrivacy(t *testing.T) {
	db := newAPIPostTestDatabase(t)
	t.Setenv("ORBIT_TOKEN_SECRET", "test-only-secret")
	token, err := GenerateToken(1, "noa")
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/posts",
		strings.NewReader(`{"content":"Hidden?","privacy":"secret"}`),
	)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})
	response := httptest.NewRecorder()

	App{DB: db}.Posts(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", response.Code, response.Body.String())
	}
}

func newAPIPostTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			dob DATE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			image_path TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			group_id INTEGER,
			privacy TEXT NOT NULL,
			like_count INTEGER NOT NULL DEFAULT 0,
			dislike_count INTEGER NOT NULL DEFAULT 0,
			comment_count INTEGER NOT NULL DEFAULT 0
		);

		CREATE TABLE follows (
			follower_id INTEGER NOT NULL,
			following_id INTEGER NOT NULL,
			PRIMARY KEY (follower_id, following_id)
		);

		CREATE TABLE profile (
			user_id INTEGER PRIMARY KEY,
			avatar_path TEXT
		);

		CREATE TABLE post_viewers (
			post_id INTEGER NOT NULL,
			viewer_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, viewer_id)
		);

		INSERT INTO users (id, email, username, first_name, last_name, dob, password)
		VALUES (1, 'noa@orbit.test', 'noa', 'Noa', 'Ferreira', '2000-01-01', 'password');
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

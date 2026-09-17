package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"social/internal/app/tokens"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestJoinRequestRollsBackWhenNotificationCreationFails(t *testing.T) {
	db := openGroupAPITestDB(t)

	_, err := db.Exec(`
		CREATE TABLE groups (id INTEGER PRIMARY KEY, creator_id INTEGER NOT NULL);
		CREATE TABLE group_members (group_id INTEGER NOT NULL, user_id INTEGER NOT NULL);
		CREATE TABLE group_join_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE UNIQUE INDEX group_join_requests_one_pending
		ON group_join_requests (group_id, user_id)
		WHERE status = 'pending';
		INSERT INTO groups (id, creator_id) VALUES (7, 1);
	`)
	if err != nil {
		t.Fatal(err)
	}

	request := authenticatedGroupRequest(t, http.MethodPost, "/api/groups/7/join", "7", 2)
	recorder := httptest.NewRecorder()

	App{DB: db}.JoinRequest(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, expected %d", recorder.Code, http.StatusInternalServerError)
	}

	var requestCount int
	if err := db.QueryRow(`
		SELECT COUNT(*)
		FROM group_join_requests
		WHERE group_id = 7 AND user_id = 2 AND status = 'pending'
	`).Scan(&requestCount); err != nil {
		t.Fatal(err)
	}
	if requestCount != 0 {
		t.Fatalf("pending request count = %d, expected 0", requestCount)
	}
}

func TestJoinRequestRejectsDuplicatePendingRequest(t *testing.T) {
	db := openGroupAPITestDB(t)

	_, err := db.Exec(`
		CREATE TABLE groups (id INTEGER PRIMARY KEY, creator_id INTEGER NOT NULL);
		CREATE TABLE group_members (group_id INTEGER NOT NULL, user_id INTEGER NOT NULL);
		CREATE TABLE group_join_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE UNIQUE INDEX group_join_requests_one_pending
		ON group_join_requests (group_id, user_id)
		WHERE status = 'pending';
		INSERT INTO groups (id, creator_id) VALUES (7, 1);
		INSERT INTO group_join_requests (group_id, user_id, status) VALUES (7, 2, 'pending');
	`)
	if err != nil {
		t.Fatal(err)
	}

	request := authenticatedGroupRequest(t, http.MethodPost, "/api/groups/7/join", "7", 2)
	recorder := httptest.NewRecorder()
	App{DB: db}.JoinRequest(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, expected %d: %s", recorder.Code, http.StatusConflict, recorder.Body.String())
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_join_requests WHERE group_id = 7 AND user_id = 2 AND status = 'pending'", 1)
}

func TestUndoJoinRequestDeletesOnlyPendingRequestAndNotification(t *testing.T) {
	db := openGroupAPITestDB(t)

	_, err := db.Exec(`
		CREATE TABLE group_join_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE UNIQUE INDEX group_join_requests_one_pending
		ON group_join_requests (group_id, user_id)
		WHERE status = 'pending';
		CREATE TABLE notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL,
			type TEXT NOT NULL,
			related_id INTEGER
		);
		INSERT INTO group_join_requests (id, group_id, user_id, status)
		VALUES (10, 7, 2, 'rejected'), (11, 7, 2, 'pending');
		INSERT INTO notifications (id, category, type, related_id)
		VALUES
			(20, 'groups', 'join_request', 10),
			(21, 'groups', 'join_request', 11);
	`)
	if err != nil {
		t.Fatal(err)
	}

	request := authenticatedGroupRequest(t, http.MethodDelete, "/api/groups/7/join", "7", 2)
	recorder := httptest.NewRecorder()

	App{DB: db}.UndoJoinRequest(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, expected %d", recorder.Code, http.StatusOK)
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM group_join_requests WHERE id = 10 AND status = 'rejected'", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_join_requests WHERE id = 11", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 20", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 21", 0)
}

func openGroupAPITestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	return db
}

func authenticatedGroupRequest(t *testing.T, method string, path string, groupID string, userID int) *http.Request {
	t.Helper()
	t.Setenv("ORBIT_TOKEN_SECRET", "group-api-test-secret")

	token, err := tokens.GenerateToken(userID)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(method, path, nil)
	request.SetPathValue("id", groupID)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})
	return request
}

func assertRowCount(t *testing.T, db *sql.DB, query string, expected int) {
	t.Helper()

	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != expected {
		t.Fatalf("row count = %d, expected %d for %q", count, expected, query)
	}
}

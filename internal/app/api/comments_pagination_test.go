package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"social/internal/app/tokens"
	"testing"
)

func TestCommentsEndpointReturnsBoundedPages(t *testing.T) {
	db := openGroupAPITestDB(t)
	_, err := db.Exec(`
		CREATE TABLE user (id INTEGER PRIMARY KEY, username TEXT, first_name TEXT, last_name TEXT);
		CREATE TABLE profile (user_id INTEGER PRIMARY KEY, avatar_path TEXT);
		CREATE TABLE posts (id INTEGER PRIMARY KEY, user_id INTEGER, group_id INTEGER, privacy TEXT);
		CREATE TABLE post_viewers (post_id INTEGER, viewer_id INTEGER);
		CREATE TABLE user_followers (follower_id INTEGER, target_id INTEGER, status INTEGER);
		CREATE TABLE comments (
			id INTEGER PRIMARY KEY,
			post_id INTEGER,
			user_id INTEGER,
			content TEXT,
			image_path TEXT,
			created_at DATETIME
		);
		INSERT INTO user (id, username, first_name, last_name) VALUES (1, 'author', 'Post', 'Author');
		INSERT INTO user (id, username, first_name, last_name) VALUES (2, 'commenter', 'Test', 'Commenter');
		INSERT INTO posts (id, user_id, privacy) VALUES (1, 1, 'public');
	`)
	if err != nil {
		t.Fatal(err)
	}

	for commentID := 1; commentID <= 45; commentID++ {
		_, err := db.Exec(`
			INSERT INTO comments (id, post_id, user_id, content, image_path, created_at)
			VALUES (?, 1, 2, ?, '', datetime('2026-01-01 00:00:00', '+' || ? || ' seconds'))
		`, commentID, fmt.Sprintf("Comment %02d", commentID), commentID)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Setenv("ORBIT_TOKEN_SECRET", "comments-pagination-test-secret")
	token, err := tokens.GenerateToken(2)
	if err != nil {
		t.Fatal(err)
	}
	app := &App{DB: db}

	readPage := func(offset int) ([]string, bool, int) {
		t.Helper()
		request := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/posts/1/comments?limit=20&offset=%d", offset), nil)
		request.SetPathValue("postID", "1")
		request.AddCookie(&http.Cookie{Name: "token", Value: token})
		request = request.WithContext(context.WithValue(request.Context(), "userID", 2))
		recorder := httptest.NewRecorder()
		app.Comments(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, expected %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}

		var response struct {
			Comments []struct {
				Content string `json:"content"`
			} `json:"comments"`
			HasMore    bool `json:"hasMore"`
			NextOffset int  `json:"nextOffset"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		contents := make([]string, 0, len(response.Comments))
		for _, comment := range response.Comments {
			contents = append(contents, comment.Content)
		}
		return contents, response.HasMore, response.NextOffset
	}

	firstPage, hasMore, nextOffset := readPage(0)
	if len(firstPage) != 20 || firstPage[0] != "Comment 01" || firstPage[19] != "Comment 20" || !hasMore || nextOffset != 20 {
		t.Fatalf("first page = (%d comments, %q..%q, hasMore=%t, nextOffset=%d)",
			len(firstPage), firstPage[0], firstPage[len(firstPage)-1], hasMore, nextOffset)
	}

	secondPage, hasMore, nextOffset := readPage(nextOffset)
	if len(secondPage) != 20 || secondPage[0] != "Comment 21" || secondPage[19] != "Comment 40" || !hasMore || nextOffset != 40 {
		t.Fatalf("second page = (%d comments, %q..%q, hasMore=%t, nextOffset=%d)",
			len(secondPage), secondPage[0], secondPage[len(secondPage)-1], hasMore, nextOffset)
	}

	lastPage, hasMore, nextOffset := readPage(nextOffset)
	if len(lastPage) != 5 || lastPage[0] != "Comment 41" || lastPage[4] != "Comment 45" || hasMore || nextOffset != 45 {
		t.Fatalf("last page = (%d comments, %q..%q, hasMore=%t, nextOffset=%d)",
			len(lastPage), lastPage[0], lastPage[len(lastPage)-1], hasMore, nextOffset)
	}
}

func TestGroupPostCommentsEndpointReturnsBoundedPages(t *testing.T) {
	db := openGroupAPITestDB(t)
	_, err := db.Exec(`
		CREATE TABLE user (id INTEGER PRIMARY KEY, username TEXT, first_name TEXT, last_name TEXT);
		CREATE TABLE profile (user_id INTEGER PRIMARY KEY, avatar_path TEXT);
		CREATE TABLE group_members (group_id INTEGER, user_id INTEGER);
		CREATE TABLE group_posts (id INTEGER PRIMARY KEY, group_id INTEGER, user_id INTEGER);
		CREATE TABLE group_post_comments (
			id INTEGER PRIMARY KEY,
			post_id INTEGER,
			user_id INTEGER,
			content TEXT,
			image_path TEXT,
			created_at DATETIME
		);
		INSERT INTO user (id, username, first_name, last_name) VALUES (1, 'owner', 'Group', 'Owner');
		INSERT INTO user (id, username, first_name, last_name) VALUES (2, 'member', 'Test', 'Member');
		INSERT INTO group_members (group_id, user_id) VALUES (1, 2);
		INSERT INTO group_posts (id, group_id, user_id) VALUES (1, 1, 1);
	`)
	if err != nil {
		t.Fatal(err)
	}

	for commentID := 1; commentID <= 45; commentID++ {
		_, err := db.Exec(`
			INSERT INTO group_post_comments (id, post_id, user_id, content, image_path, created_at)
			VALUES (?, 1, 2, ?, '', datetime('2026-01-01 00:00:00', '+' || ? || ' seconds'))
		`, commentID, fmt.Sprintf("Group comment %02d", commentID), commentID)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Setenv("ORBIT_TOKEN_SECRET", "group-comments-pagination-test-secret")
	token, err := tokens.GenerateToken(2)
	if err != nil {
		t.Fatal(err)
	}
	app := &App{DB: db}

	readPage := func(offset int) ([]string, bool, int) {
		t.Helper()
		request := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/api/groups/1/posts/1/comments?limit=20&offset=%d", offset), nil)
		request.SetPathValue("id", "1")
		request.SetPathValue("postID", "1")
		request.AddCookie(&http.Cookie{Name: "token", Value: token})
		request = request.WithContext(context.WithValue(request.Context(), "userID", 2))
		recorder := httptest.NewRecorder()
		app.GetGroupPostComments(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, expected %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}

		var response struct {
			Comments []struct {
				Content string `json:"content"`
			} `json:"comments"`
			HasMore    bool `json:"hasMore"`
			NextOffset int  `json:"nextOffset"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		contents := make([]string, 0, len(response.Comments))
		for _, comment := range response.Comments {
			contents = append(contents, comment.Content)
		}
		return contents, response.HasMore, response.NextOffset
	}

	firstPage, hasMore, nextOffset := readPage(0)
	if len(firstPage) != 20 || firstPage[0] != "Group comment 01" || firstPage[19] != "Group comment 20" || !hasMore || nextOffset != 20 {
		t.Fatalf("first group page = (%d comments, %q..%q, hasMore=%t, nextOffset=%d)",
			len(firstPage), firstPage[0], firstPage[len(firstPage)-1], hasMore, nextOffset)
	}

	secondPage, hasMore, nextOffset := readPage(nextOffset)
	if len(secondPage) != 20 || secondPage[0] != "Group comment 21" || secondPage[19] != "Group comment 40" || !hasMore || nextOffset != 40 {
		t.Fatalf("second group page = (%d comments, %q..%q, hasMore=%t, nextOffset=%d)",
			len(secondPage), secondPage[0], secondPage[len(secondPage)-1], hasMore, nextOffset)
	}

	lastPage, hasMore, nextOffset := readPage(nextOffset)
	if len(lastPage) != 5 || lastPage[0] != "Group comment 41" || lastPage[4] != "Group comment 45" || hasMore || nextOffset != 45 {
		t.Fatalf("last group page = (%d comments, %q..%q, hasMore=%t, nextOffset=%d)",
			len(lastPage), lastPage[0], lastPage[len(lastPage)-1], hasMore, nextOffset)
	}
}

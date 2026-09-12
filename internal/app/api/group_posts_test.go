package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGroupPostsRequireMembershipAndSupportTextOrImage(t *testing.T) {
	db := openGroupPostsAPITestDB(t)
	app := App{DB: db}

	recorder := httptest.NewRecorder()
	app.CreateGroupPost(recorder, groupPostRequest(t, http.MethodPost, "7", "", 1, "First post", nil))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("text post status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	app.CreateGroupPost(recorder, groupPostRequest(t, http.MethodPost, "7", "", 1, "Second post", nil))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("second post status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	app.GetGroupPosts(recorder, groupPostRequest(t, http.MethodGet, "7", "", 1, "", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, expected %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	var listed struct {
		Posts []struct {
			ID      int64  `json:"id"`
			Content string `json:"content"`
			IsOwner bool   `json:"isOwner"`
		} `json:"posts"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&listed); err != nil {
		t.Fatal(err)
	}
	if len(listed.Posts) != 2 || listed.Posts[0].Content != "Second post" || listed.Posts[1].Content != "First post" {
		t.Fatalf("posts were not newest first: %+v", listed.Posts)
	}
	if !listed.Posts[0].IsOwner || !listed.Posts[1].IsOwner {
		t.Fatalf("the post author was not marked as owner: %+v", listed.Posts)
	}

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		recorder = httptest.NewRecorder()
		request := groupPostRequest(t, method, "7", "", 3, "not allowed", nil)
		if method == http.MethodGet {
			app.GetGroupPosts(recorder, request)
		} else {
			app.CreateGroupPost(recorder, request)
		}
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("non-member %s status = %d, expected %d", method, recorder.Code, http.StatusForbidden)
		}
	}

	recorder = httptest.NewRecorder()
	app.CreateGroupPost(recorder, groupPostRequest(t, http.MethodPost, "7", "", 1, "", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("empty post status = %d, expected %d", recorder.Code, http.StatusBadRequest)
	}

	recorder = httptest.NewRecorder()
	app.CreateGroupPost(recorder, groupPostRequest(t, http.MethodPost, "7", "", 1, "", minimalGIF()))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("image-only post status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	var created struct {
		Post struct {
			ImagePath string `json:"imagePath"`
		} `json:"post"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	cleanupGroupPostTestUpload(t, created.Post.ImagePath)
}

func TestGroupPostCommentsAreMemberOnlyAndGroupScoped(t *testing.T) {
	db := openGroupPostsAPITestDB(t)
	app := App{DB: db}

	result, err := db.Exec(`INSERT INTO group_posts (group_id, user_id, content) VALUES (7, 1, 'Group seven post')`)
	if err != nil {
		t.Fatal(err)
	}
	postID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		recorder := httptest.NewRecorder()
		request := groupPostRequest(t, method, "7", strconv.FormatInt(postID, 10), 3, "not allowed", nil)
		if method == http.MethodGet {
			app.GetGroupPostComments(recorder, request)
		} else {
			app.CreateGroupPostComment(recorder, request)
		}
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("non-member comment %s status = %d, expected %d", method, recorder.Code, http.StatusForbidden)
		}
	}

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		recorder := httptest.NewRecorder()
		request := groupPostRequest(t, method, "8", strconv.FormatInt(postID, 10), 2, "wrong group", nil)
		if method == http.MethodGet {
			app.GetGroupPostComments(recorder, request)
		} else {
			app.CreateGroupPostComment(recorder, request)
		}
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("wrong-group comment %s status = %d, expected %d", method, recorder.Code, http.StatusNotFound)
		}
	}

	recorder := httptest.NewRecorder()
	app.CreateGroupPostComment(recorder, groupPostRequest(t, http.MethodPost, "7", strconv.FormatInt(postID, 10), 2, "First comment", nil))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("comment status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	recorder = httptest.NewRecorder()
	app.CreateGroupPostComment(recorder, groupPostRequest(t, http.MethodPost, "7", strconv.FormatInt(postID, 10), 1, "Second comment", nil))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("second comment status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	app.GetGroupPostComments(recorder, groupPostRequest(t, http.MethodGet, "7", strconv.FormatInt(postID, 10), 1, "", nil))
	var listed struct {
		Comments []struct {
			Content string `json:"content"`
			IsOwner bool   `json:"isOwner"`
		} `json:"comments"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&listed); err != nil {
		t.Fatal(err)
	}
	if len(listed.Comments) != 2 || listed.Comments[0].Content != "First comment" || listed.Comments[1].Content != "Second comment" {
		t.Fatalf("comments were not oldest first: %+v", listed.Comments)
	}
	if listed.Comments[0].IsOwner || !listed.Comments[1].IsOwner {
		t.Fatalf("comment ownership flags were incorrect: %+v", listed.Comments)
	}

	recorder = httptest.NewRecorder()
	app.CreateGroupPostComment(recorder, groupPostRequest(t, http.MethodPost, "7", strconv.FormatInt(postID, 10), 1, "", minimalGIF()))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("image-only comment status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	var created struct {
		Comment struct {
			ImagePath string `json:"imagePath"`
		} `json:"comment"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	cleanupGroupPostTestUpload(t, created.Comment.ImagePath)

	recorder = httptest.NewRecorder()
	app.CreateGroupPostComment(recorder, groupPostRequest(t, http.MethodPost, "7", strconv.FormatInt(postID, 10), 1, "", nil))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("empty comment status = %d, expected %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGroupPostListsDoNotLeakAcrossGroups(t *testing.T) {
	db := openGroupPostsAPITestDB(t)
	app := App{DB: db}
	if _, err := db.Exec(`
		INSERT INTO group_posts (group_id, user_id, content)
		VALUES (7, 1, 'Group seven'), (8, 2, 'Group eight')
	`); err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	app.GetGroupPosts(recorder, groupPostRequest(t, http.MethodGet, "8", "", 2, "", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, expected %d", recorder.Code, http.StatusOK)
	}
	var response struct {
		Posts []struct {
			Content string `json:"content"`
		} `json:"posts"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if len(response.Posts) != 1 || response.Posts[0].Content != "Group eight" {
		t.Fatalf("group list leaked posts: %+v", response.Posts)
	}
}

func TestDeleteGroupPostEnforcesMembershipOwnershipAndGroup(t *testing.T) {
	t.Run("author deletes post and cascading comments", func(t *testing.T) {
		db := openGroupPostsAPITestDB(t)
		seedGroupPostDeletionData(t, db)

		recorder := httptest.NewRecorder()
		App{DB: db}.DeleteGroupPost(recorder, groupPostDeleteRequest(t, "7", "10", "", 1))
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_posts WHERE id = 10", 0)
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_post_comments WHERE post_id = 10", 0)
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_posts WHERE id = 20 AND group_id = 8", 1)
	})

	tests := []struct {
		name     string
		groupID  string
		postID   string
		userID   int
		expected int
	}{
		{name: "another group member", groupID: "7", postID: "10", userID: 2, expected: http.StatusForbidden},
		{name: "non-member", groupID: "7", postID: "10", userID: 3, expected: http.StatusForbidden},
		{name: "wrong group", groupID: "8", postID: "10", userID: 3, expected: http.StatusNotFound},
		{name: "missing post", groupID: "7", postID: "999", userID: 1, expected: http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := openGroupPostsAPITestDB(t)
			seedGroupPostDeletionData(t, db)

			recorder := httptest.NewRecorder()
			App{DB: db}.DeleteGroupPost(recorder, groupPostDeleteRequest(t, test.groupID, test.postID, "", test.userID))
			if recorder.Code != test.expected {
				t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, test.expected, recorder.Body.String())
			}
			assertRowCount(t, db, "SELECT COUNT(*) FROM group_posts WHERE id = 10", 1)
			assertRowCount(t, db, "SELECT COUNT(*) FROM group_posts WHERE id = 20", 1)
		})
	}
}

func TestDeleteGroupPostCommentEnforcesExactOwnershipAndRelationships(t *testing.T) {
	t.Run("comment author deletes only their comment", func(t *testing.T) {
		db := openGroupPostsAPITestDB(t)
		seedGroupPostDeletionData(t, db)
		app := App{DB: db}

		recorder := httptest.NewRecorder()
		app.DeleteGroupPostComment(recorder, groupPostDeleteRequest(t, "7", "10", "100", 1))
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_post_comments WHERE id = 100", 0)
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_post_comments WHERE id = 101", 1)

		recorder = httptest.NewRecorder()
		app.GetGroupPosts(recorder, groupPostRequest(t, http.MethodGet, "7", "", 1, "", nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("post refresh status = %d, expected %d", recorder.Code, http.StatusOK)
		}
		var response struct {
			Posts []struct {
				ID           int64 `json:"id"`
				CommentCount int   `json:"commentCount"`
			} `json:"posts"`
		}
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		for _, post := range response.Posts {
			if post.ID == 10 && post.CommentCount != 1 {
				t.Fatalf("comment count = %d, expected 1", post.CommentCount)
			}
		}
	})

	tests := []struct {
		name      string
		groupID   string
		postID    string
		commentID string
		userID    int
		expected  int
	}{
		{name: "another member", groupID: "7", postID: "10", commentID: "100", userID: 2, expected: http.StatusForbidden},
		{name: "post author is not comment owner", groupID: "7", postID: "10", commentID: "101", userID: 1, expected: http.StatusForbidden},
		{name: "non-member", groupID: "7", postID: "10", commentID: "100", userID: 3, expected: http.StatusForbidden},
		{name: "wrong group", groupID: "8", postID: "20", commentID: "100", userID: 3, expected: http.StatusNotFound},
		{name: "wrong post", groupID: "7", postID: "11", commentID: "100", userID: 2, expected: http.StatusNotFound},
		{name: "comment from another group", groupID: "7", postID: "10", commentID: "200", userID: 1, expected: http.StatusNotFound},
		{name: "missing comment", groupID: "7", postID: "10", commentID: "999", userID: 1, expected: http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := openGroupPostsAPITestDB(t)
			seedGroupPostDeletionData(t, db)

			recorder := httptest.NewRecorder()
			App{DB: db}.DeleteGroupPostComment(recorder, groupPostDeleteRequest(t, test.groupID, test.postID, test.commentID, test.userID))
			if recorder.Code != test.expected {
				t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, test.expected, recorder.Body.String())
			}
			assertRowCount(t, db, "SELECT COUNT(*) FROM group_post_comments WHERE id IN (100, 101, 200)", 3)
		})
	}
}

func openGroupPostsAPITestDB(t *testing.T) *sql.DB {
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
		CREATE TABLE groups (id INTEGER PRIMARY KEY, creator_id INTEGER NOT NULL);
		CREATE TABLE group_members (group_id INTEGER NOT NULL, user_id INTEGER NOT NULL, PRIMARY KEY (group_id, user_id));
		CREATE TABLE group_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			image_path TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		);
		CREATE TABLE group_post_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			image_path TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES group_posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
		);
		INSERT INTO user (id, username, first_name, last_name) VALUES
			(1, 'owner', 'Group', 'Owner'),
			(2, 'member', 'Group', 'Member'),
			(3, 'outsider', 'Outside', 'User');
		INSERT INTO profile (user_id, avatar_path) VALUES
			(1, 'avatars/default.png'),
			(2, 'avatars/default.png'),
			(3, 'avatars/default.png');
		INSERT INTO groups (id, creator_id) VALUES (7, 1), (8, 2);
		INSERT INTO group_members (group_id, user_id) VALUES (7, 1), (7, 2), (8, 2), (8, 3);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func seedGroupPostDeletionData(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		INSERT INTO group_posts (id, group_id, user_id, content) VALUES
			(10, 7, 1, 'Owner post'),
			(11, 7, 2, 'Member post'),
			(20, 8, 3, 'Other group post');
		INSERT INTO group_post_comments (id, post_id, user_id, content) VALUES
			(100, 10, 1, 'Owner comment'),
			(101, 10, 2, 'Other member comment'),
			(110, 11, 2, 'Second post comment'),
			(200, 20, 3, 'Other group comment');
	`)
	if err != nil {
		t.Fatal(err)
	}
}

func groupPostRequest(t *testing.T, method, groupID, postID string, userID int, content string, image []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("content", content); err != nil {
		t.Fatal(err)
	}
	if image != nil {
		part, err := writer.CreateFormFile("image", "test.gif")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(image); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(method, "/api/groups/"+groupID+"/posts", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.SetPathValue("id", groupID)
	if postID != "" {
		request.SetPathValue("postID", postID)
	}
	ctx := context.WithValue(request.Context(), "userID", userID)
	return request.WithContext(ctx)
}

func groupPostDeleteRequest(t *testing.T, groupID, postID, commentID string, userID int) *http.Request {
	t.Helper()
	request := groupPostRequest(t, http.MethodDelete, groupID, postID, userID, "", nil)
	if commentID != "" {
		request.SetPathValue("commentID", commentID)
	}
	return request
}

func minimalGIF() []byte {
	return []byte("GIF89a\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\xff\xff\xff!\xf9\x04\x01\x00\x00\x00\x00,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02D\x01\x00;")
}

func cleanupGroupPostTestUpload(t *testing.T, imagePath string) {
	t.Helper()
	if imagePath == "" {
		t.Fatal("created image path was empty")
	}
	path := filepath.Join("uploads", filepath.FromSlash(imagePath))
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove test upload: %v", err)
	}
	_ = os.Remove(filepath.Dir(path))
	_ = os.Remove("uploads")
}

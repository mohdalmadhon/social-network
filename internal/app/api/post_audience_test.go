package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social/internal/app/tokens"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestFollowersEndpointPaginatesTheActualAudienceAccounts(t *testing.T) {
	db := openGroupAPITestDB(t)
	_, err := db.Exec(`
		CREATE TABLE user (id INTEGER PRIMARY KEY, first_name TEXT, last_name TEXT);
		CREATE TABLE profile (user_id INTEGER PRIMARY KEY, avatar_path TEXT);
		CREATE TABLE user_followers (follower_id INTEGER, target_id INTEGER, status INTEGER);
		INSERT INTO user (id, first_name, last_name) VALUES (1, 'Post', 'Author');
	`)
	if err != nil {
		t.Fatal(err)
	}

	for followerID := 101; followerID <= 130; followerID++ {
		if _, err := db.Exec(`INSERT INTO user (id, first_name, last_name) VALUES (?, 'Follower', 'Account')`, followerID); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`INSERT INTO user_followers (follower_id, target_id, status) VALUES (?, 1, 1)`, followerID); err != nil {
			t.Fatal(err)
		}
	}

	t.Setenv("ORBIT_TOKEN_SECRET", "post-audience-test-secret")
	token, err := tokens.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}
	app := &App{DB: db}

	readPage := func(offset int) []int {
		t.Helper()
		request := httptest.NewRequest(http.MethodGet, "/api/profile/follow?count=25&offset="+strconv.Itoa(offset), nil)
		request.AddCookie(&http.Cookie{Name: "token", Value: token})
		request = request.WithContext(context.WithValue(request.Context(), "userID", 1))
		recorder := httptest.NewRecorder()
		app.GetFollowers(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, expected %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}

		var response struct {
			Data []struct {
				ID int
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		ids := make([]int, 0, len(response.Data))
		for _, follower := range response.Data {
			ids = append(ids, follower.ID)
		}
		return ids
	}

	firstPage := readPage(0)
	if len(firstPage) != 25 || firstPage[0] != 101 || firstPage[24] != 125 {
		t.Fatalf("first follower page = %v, expected account IDs 101 through 125", firstPage)
	}

	secondPage := readPage(25)
	if len(secondPage) != 5 || secondPage[0] != 126 || secondPage[4] != 130 {
		t.Fatalf("second follower page = %v, expected account IDs 126 through 130", secondPage)
	}
}

package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// Use the real HTTP routes, registration, password hashing, cookies and migrations.
// Nothing in this test touches the developer's database or existing accounts.
func TestTwoAccountFollowNotifications(t *testing.T) {
	t.Setenv("ORBIT_TOKEN_SECRET", "isolated-integration-test-secret")
	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	files, err := filepath.Glob("../../../internal/migrations/*.up.sql")
	if err != nil || len(files) == 0 {
		t.Fatal("migrations not found", err)
	}
	for _, name := range files {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = db.Exec(string(data)); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	server := httptest.NewServer(StartServer(db))
	defer server.Close()
	register := func(name string) *http.Client {
		t.Helper()
		var body bytes.Buffer
		form := multipart.NewWriter(&body)
		for key, value := range map[string]string{"FirstName": name, "LastName": "Tester", "UserName": name, "Email": name + "@example.test", "Password": "TestPass123!", "dob": "2000-01-01"} {
			if err := form.WriteField(key, value); err != nil {
				t.Fatal(err)
			}
		}
		form.Close()
		response, err := http.Post(server.URL+"/api/user", form.FormDataContentType(), &body)
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(response.Body)
		response.Body.Close()
		if response.StatusCode != 201 {
			t.Fatalf("registration: %d %s", response.StatusCode, data)
		}
		jar, _ := cookiejar.New(nil)
		client := &http.Client{Jar: jar}
		login := bytes.NewBufferString(fmt.Sprintf(`{"Identifier":%q,"Pass":"TestPass123!"}`, name))
		response, err = client.Post(server.URL+"/api/session", "application/json", login)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != 200 {
			t.Fatalf("login: %d", response.StatusCode)
		}
		return client
	}
	alice, bob := register("alice"), register("bob")
	if _, err := db.Exec(`UPDATE profile SET is_private = 1 WHERE user_id = 2`); err != nil {
		t.Fatal(err)
	}
	send := func(client *http.Client, method, path, body string, expected int) map[string]any {
		t.Helper()
		request, _ := http.NewRequest(method, server.URL+path, bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		data, _ := io.ReadAll(response.Body)
		if response.StatusCode != expected {
			t.Fatalf("%s %s: %d %s", method, path, response.StatusCode, data)
		}
		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatal(err)
		}
		return result
	}
	send(alice, "POST", "/api/profile/follow?targetid=2", "", 200)
	send(alice, "POST", "/api/profile/follow?targetid=2", "", 200)
	result := send(bob, "GET", "/api/notifications", "", 200)
	items := result["notifications"].([]any)
	if len(items) != 1 || result["unreadCount"] != float64(1) {
		t.Fatalf("expected one notification: %+v", result)
	}
	id := int(items[0].(map[string]any)["id"].(float64))
	actionPath := fmt.Sprintf("/api/notifications/%d/action", id)
	send(alice, "PATCH", actionPath, `{"action":"accept"}`, 404)
	send(bob, "PATCH", actionPath, `{"action":"accept"}`, 200)
	send(bob, "PATCH", actionPath, `{"action":"accept"}`, 404)
	result = send(bob, "GET", "/api/notifications", "", 200)
	if result["unreadCount"] != float64(0) {
		t.Fatal("accepted notification is still unread")
	}
	item := result["notifications"].([]any)[0].(map[string]any)
	if item["followStatus"] != float64(1) {
		t.Fatal("acceptance did not persist")
	}
	send(alice, "POST", "/api/profile/follow?targetid=2", "", 200)
	var status int
	if err := db.QueryRow(`SELECT status FROM user_followers WHERE follower_id=1 AND target_id=2`).Scan(&status); err != nil || status != 1 {
		t.Fatal("duplicate follow downgraded acceptance", err)
	}
	send(alice, "DELETE", "/api/profile/follow?targetid=2", "", 200)
	send(bob, "PATCH", actionPath, `{"action":"accept"}`, 404)
	send(alice, "POST", "/api/profile/follow?targetid=2", "", 200)
	result = send(bob, "GET", "/api/notifications", "", 200)
	id = int(result["notifications"].([]any)[0].(map[string]any)["id"].(float64))
	send(bob, "PATCH", fmt.Sprintf("/api/notifications/%d/action", id), `{"action":"decline"}`, 200)
	result = send(bob, "GET", "/api/notifications", "", 200)
	if len(result["notifications"].([]any)) != 0 {
		t.Fatal("declined request remains actionable")
	}
	// Invitations and event notifications use the same two authenticated clients.
	send(alice, "POST", "/api/groups", `{"title":"Test orbit","description":"Isolated test group"}`, 201)
	send(bob, "POST", "/api/groups/1/invitations", `{"userId":1}`, 403)
	send(alice, "POST", "/api/groups/1/invitations", `{"userId":2}`, 201)
	result = send(bob, "GET", "/api/notifications", "", 200)
	id = int(result["notifications"].([]any)[0].(map[string]any)["id"].(float64))
	send(bob, "PATCH", fmt.Sprintf("/api/notifications/%d/action", id), `{"action":"join"}`, 200)
	send(alice, "POST", "/api/groups/1/events", `{"title":"Meetup","description":"Two account event","startsAt":"2099-01-01T12:00:00Z"}`, 201)
	result = send(bob, "GET", "/api/notifications", "", 200)
	id = int(result["notifications"].([]any)[0].(map[string]any)["id"].(float64))
	send(bob, "PATCH", fmt.Sprintf("/api/notifications/%d/action", id), `{"action":"rsvp"}`, 200)
	result = send(bob, "GET", "/api/groups/1/events", "", 200)
	if result["events"].([]any)[0].(map[string]any)["response"] != "going" {
		t.Fatal("RSVP did not persist")
	}
	send(alice, "POST", "/api/groups", `{"title":"Request group","description":"Join request test"}`, 201)
	send(bob, "GET", "/api/groups/2/events", "", 403)
	send(bob, "POST", "/api/groups/2/join-requests", "", 201)
	result = send(alice, "GET", "/api/notifications", "", 200)
	items = result["notifications"].([]any)
	id = int(items[0].(map[string]any)["id"].(float64))
	send(alice, "PATCH", fmt.Sprintf("/api/notifications/%d/action", id), `{"action":"accept"}`, 200)
	send(bob, "GET", "/api/groups/2/events", "", 200)
}

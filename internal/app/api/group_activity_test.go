package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInviteUsersRequiresMembershipAndReportsOnlyPendingInvitations(t *testing.T) {
	db := newGroupActivityTestDB(t)
	app := App{DB: db}

	response := httptest.NewRecorder()
	app.GetInviteUsers(response, groupActivityRequest(http.MethodGet, "7", "", 2, ""))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, expected %d: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var result struct {
		Users []struct {
			ID           int    `json:"id"`
			IsInvited    bool   `json:"isInvited"`
			InvitationID *int64 `json:"invitationId"`
		} `json:"users"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	users := make(map[int]struct {
		invited      bool
		invitationID *int64
	})
	for _, user := range result.Users {
		users[user.ID] = struct {
			invited      bool
			invitationID *int64
		}{user.IsInvited, user.InvitationID}
	}
	if _, found := users[1]; found {
		t.Fatal("group owner appeared in invite users")
	}
	if _, found := users[2]; found {
		t.Fatal("existing group member appeared in invite users")
	}
	for _, userID := range []int{3, 4} {
		if user, found := users[userID]; !found || user.invited || user.invitationID != nil {
			t.Fatalf("historical invitation marked user %d as pending: %+v", userID, user)
		}
	}
	if user := users[5]; !user.invited || user.invitationID == nil || *user.invitationID != 12 {
		t.Fatalf("pending invitation was not reported correctly: %+v", user)
	}

	response = httptest.NewRecorder()
	app.GetInviteUsers(response, groupActivityRequest(http.MethodGet, "7", "", 6, ""))
	if response.Code != http.StatusForbidden {
		t.Fatalf("nonmember status = %d, expected %d", response.Code, http.StatusForbidden)
	}

	if _, err := db.Exec(`INSERT INTO group_members (group_id, user_id) VALUES (7, 5)`); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	app.GetInviteUsers(response, groupActivityRequest(http.MethodGet, "7", "", 2, ""))
	if strings.Contains(response.Body.String(), `"id":5`) {
		t.Fatal("newly accepted member remained in invite users")
	}
}

func TestInvitationCreateUndoAndRollbackUseExactPendingInvitation(t *testing.T) {
	db := newGroupActivityTestDB(t)
	app := App{DB: db}

	response := httptest.NewRecorder()
	app.InviteGroupMember(response, groupActivityRequest(http.MethodPost, "7", "", 2, `{"userId":3}`))
	if response.Code != http.StatusCreated {
		t.Fatalf("invite status = %d, expected %d: %s", response.Code, http.StatusCreated, response.Body.String())
	}
	var sent struct {
		InvitationID int64 `json:"invitationId"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &sent); err != nil {
		t.Fatal(err)
	}
	if sent.InvitationID <= 12 {
		t.Fatalf("invitation id = %d, expected a new history row", sent.InvitationID)
	}
	assertRowCount(t, db, fmt.Sprintf("SELECT COUNT(*) FROM group_invitations WHERE id = %d AND status = 'pending'", sent.InvitationID), 1)
	assertRowCount(t, db, fmt.Sprintf("SELECT COUNT(*) FROM notifications WHERE type = 'invitation' AND related_id = %d", sent.InvitationID), 1)

	response = httptest.NewRecorder()
	app.InviteGroupMember(response, groupActivityRequest(http.MethodPost, "7", "", 2, `{"userId":3}`))
	if response.Code != http.StatusConflict {
		t.Fatalf("duplicate pending status = %d, expected %d", response.Code, http.StatusConflict)
	}

	response = httptest.NewRecorder()
	app.UndoInvitation(response, groupActivityRequest(http.MethodDelete, "7", "12", 2, ""))
	if response.Code != http.StatusOK {
		t.Fatalf("undo status = %d, expected %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 12", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 21", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id IN (10, 11)", 2)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 20", 1)

	response = httptest.NewRecorder()
	app.InviteGroupMember(response, groupActivityRequest(http.MethodPost, "7", "", 2, `{"userId":5}`))
	if response.Code != http.StatusCreated {
		t.Fatalf("reinvite after undo status = %d, expected %d: %s", response.Code, http.StatusCreated, response.Body.String())
	}

	response = httptest.NewRecorder()
	app.UndoInvitation(response, groupActivityRequest(http.MethodDelete, "7", fmt.Sprint(sent.InvitationID), 6, ""))
	if response.Code != http.StatusForbidden {
		t.Fatalf("nonmember undo status = %d, expected %d", response.Code, http.StatusForbidden)
	}
	if _, err := db.Exec(`INSERT INTO group_members (group_id, user_id) VALUES (8, 2)`); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	app.UndoInvitation(response, groupActivityRequest(http.MethodDelete, "8", fmt.Sprint(sent.InvitationID), 2, ""))
	if response.Code != http.StatusNotFound {
		t.Fatalf("wrong-group undo status = %d, expected %d", response.Code, http.StatusNotFound)
	}

	if _, err := db.Exec(`DROP TABLE notifications`); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	app.InviteGroupMember(response, groupActivityRequest(http.MethodPost, "7", "", 2, `{"userId":4}`))
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("failed-notification status = %d, expected %d", response.Code, http.StatusInternalServerError)
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE group_id = 7 AND user_id = 4 AND status = 'pending'", 0)
}

func TestGroupEventsEnforceMembershipAndEventOwnership(t *testing.T) {
	db := newGroupActivityTestDB(t)
	app := App{DB: db}

	response := httptest.NewRecorder()
	app.GroupEvents(response, groupActivityRequest(http.MethodGet, "7", "", 2, ""))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"id":30`) || strings.Contains(response.Body.String(), `"id":31`) {
		t.Fatalf("group event list leaked or omitted an event: %d %s", response.Code, response.Body.String())
	}
	var listed struct {
		Events []groupEvent `json:"events"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	if len(listed.Events) != 1 || listed.Events[0].GoingCount != 1 || listed.Events[0].NotGoingCount != 1 {
		t.Fatalf("event counts were not returned: %+v", listed.Events)
	}
	if listed.Events[0].Response != "declined" || listed.Events[0].IsCreator {
		t.Fatalf("member event state was incorrect: %+v", listed.Events[0])
	}
	if len(listed.Events[0].GoingUsers) != 1 || listed.Events[0].GoingUsers[0].ID != 1 || len(listed.Events[0].NotGoingUsers) != 1 || listed.Events[0].NotGoingUsers[0].ID != 2 {
		t.Fatalf("event voter lists were incorrect: %+v", listed.Events[0])
	}

	response = httptest.NewRecorder()
	app.GroupEvents(response, groupActivityRequest(http.MethodGet, "7", "", 6, ""))
	if response.Code != http.StatusForbidden {
		t.Fatalf("nonmember event list status = %d, expected %d", response.Code, http.StatusForbidden)
	}
	response = httptest.NewRecorder()
	app.GroupEvents(response, groupActivityRequest(http.MethodPost, "7", "", 6, `{"title":"Blocked","description":"Not allowed","startsAt":"2099-01-03T12:00:00Z"}`))
	if response.Code != http.StatusForbidden {
		t.Fatalf("nonmember event create status = %d, expected %d", response.Code, http.StatusForbidden)
	}

	response = httptest.NewRecorder()
	app.EventRSVP(response, groupActivityRequest(http.MethodPatch, "", "30", 6, `{"response":"going"}`))
	if response.Code != http.StatusNotFound {
		t.Fatalf("cross-group RSVP status = %d, expected %d", response.Code, http.StatusNotFound)
	}
	response = httptest.NewRecorder()
	app.EventRSVP(response, groupActivityRequest(http.MethodPatch, "", "31", 2, `{"response":"going"}`))
	if response.Code != http.StatusNotFound {
		t.Fatalf("other-group RSVP status = %d, expected %d", response.Code, http.StatusNotFound)
	}
	response = httptest.NewRecorder()
	app.EventRSVP(response, groupActivityRequest(http.MethodPatch, "", "30", 2, `{"response":"going"}`))
	if response.Code != http.StatusOK {
		t.Fatalf("member RSVP status = %d, expected %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM event_rsvps WHERE event_id = 30 AND user_id = 2", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM event_rsvps WHERE event_id = 30 AND user_id = 2 AND response = 'going'", 1)

	response = httptest.NewRecorder()
	app.RemoveEventRSVP(response, groupActivityRequest(http.MethodDelete, "7", "30", 2, ""))
	if response.Code != http.StatusOK {
		t.Fatalf("remove RSVP status = %d, expected %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM event_rsvps WHERE event_id = 30 AND user_id = 2", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM event_rsvps WHERE event_id = 30 AND user_id = 1", 1)

	response = httptest.NewRecorder()
	app.RemoveEventRSVP(response, groupActivityRequest(http.MethodDelete, "8", "30", 6, ""))
	if response.Code != http.StatusNotFound {
		t.Fatalf("mismatched RSVP removal status = %d, expected %d", response.Code, http.StatusNotFound)
	}
}

func TestDeleteEventRequiresExactGroupAndCreatorAndCleansRelatedData(t *testing.T) {
	db := newGroupActivityTestDB(t)
	app := App{DB: db}

	response := httptest.NewRecorder()
	app.DeleteEvent(response, groupActivityRequest(http.MethodDelete, "7", "30", 2, ""))
	if response.Code != http.StatusForbidden {
		t.Fatalf("member delete status = %d, expected %d", response.Code, http.StatusForbidden)
	}
	response = httptest.NewRecorder()
	app.DeleteEvent(response, groupActivityRequest(http.MethodDelete, "8", "30", 6, ""))
	if response.Code != http.StatusNotFound {
		t.Fatalf("mismatched group delete status = %d, expected %d", response.Code, http.StatusNotFound)
	}
	response = httptest.NewRecorder()
	app.DeleteEvent(response, groupActivityRequest(http.MethodDelete, "7", "30", 1, ""))
	if response.Code != http.StatusOK {
		t.Fatalf("creator delete status = %d, expected %d: %s", response.Code, http.StatusOK, response.Body.String())
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM events WHERE id = 30", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM event_rsvps WHERE event_id = 30", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 22", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 23", 1)
}

func newGroupActivityTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE user (id INTEGER PRIMARY KEY, username TEXT, first_name TEXT NOT NULL, last_name TEXT NOT NULL);
		CREATE TABLE profile (user_id INTEGER PRIMARY KEY, avatar_path TEXT);
		CREATE TABLE groups (id INTEGER PRIMARY KEY, title TEXT NOT NULL, creator_id INTEGER NOT NULL);
		CREATE TABLE group_members (group_id INTEGER NOT NULL, user_id INTEGER NOT NULL, PRIMARY KEY (group_id, user_id));
		CREATE TABLE group_invitations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			inviter_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE UNIQUE INDEX group_invitations_one_pending ON group_invitations(group_id, user_id) WHERE status = 'pending';
		CREATE TABLE notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			actor_id INTEGER,
			category TEXT NOT NULL,
			type TEXT NOT NULL,
			message TEXT NOT NULL,
			related_id INTEGER
		);
		CREATE TRIGGER delete_group_invitation_notifications AFTER DELETE ON group_invitations BEGIN
			DELETE FROM notifications WHERE category = 'groups' AND type = 'invitation' AND related_id = OLD.id;
		END;
		CREATE TABLE events (id INTEGER PRIMARY KEY, group_id INTEGER NOT NULL, creator_id INTEGER NOT NULL, title TEXT NOT NULL, content TEXT NOT NULL, starts_at TEXT);
		CREATE TABLE event_rsvps (
			event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL,
			response TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (event_id, user_id)
		);
		CREATE TRIGGER delete_event_notifications AFTER DELETE ON events BEGIN
			DELETE FROM notifications WHERE category = 'events' AND related_id = OLD.id;
		END;
		INSERT INTO user (id, username, first_name, last_name) VALUES
			(1, 'owner', 'Group', 'Owner'),
			(2, 'member', 'Group', 'Member'),
			(3, 'declined', 'Declined', 'Invite'),
			(4, 'accepted', 'Accepted', 'Invite'),
			(5, 'pending', 'Pending', 'Invite'),
			(6, 'other', 'Other', 'Member');
		INSERT INTO profile (user_id, avatar_path) SELECT id, '' FROM user;
		INSERT INTO groups (id, title, creator_id) VALUES (7, 'Main group', 1), (8, 'Other group', 6);
		INSERT INTO group_members (group_id, user_id) VALUES (7, 1), (7, 2), (8, 6);
		INSERT INTO group_invitations (id, group_id, user_id, inviter_id, status) VALUES
			(10, 7, 3, 1, 'declined'),
			(11, 7, 4, 1, 'accepted'),
			(12, 7, 5, 1, 'pending');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id) VALUES
			(20, 3, 1, 'groups', 'invitation', 'Old invitation', 10),
			(21, 5, 1, 'groups', 'invitation', 'Pending invitation', 12);
		INSERT INTO events (id, group_id, creator_id, title, content, starts_at) VALUES
			(30, 7, 1, 'Main event', 'Main group event', '2099-01-01T12:00:00Z'),
			(31, 8, 6, 'Other event', 'Other group event', '2099-01-02T12:00:00Z');
		INSERT INTO event_rsvps (event_id, user_id, response) VALUES
			(30, 1, 'going'),
			(30, 2, 'declined'),
			(31, 6, 'going');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id) VALUES
			(22, 2, 1, 'events', 'event_created', 'Main event', 30),
			(23, 6, 6, 'events', 'event_created', 'Other event', 31);
	`)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func groupActivityRequest(method, groupID, resourceID string, userID int, body string) *http.Request {
	request := httptest.NewRequest(method, "/", strings.NewReader(body))
	if groupID != "" {
		request.SetPathValue("id", groupID)
	}
	if resourceID != "" {
		request.SetPathValue("invitationID", resourceID)
		request.SetPathValue("eventID", resourceID)
	}
	return request.WithContext(context.WithValue(request.Context(), "userID", userID))
}

package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social/internal/app/tokens"
	"strconv"
	"strings"
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

func TestUserInviteCreatesInvitationAndNotificationForMember(t *testing.T) {
	db := openInvitationAPITestDB(t, true)
	app := App{DB: db}

	request := groupInvitationRequest(http.MethodPost, "/api/groups/7/invitations", "7", 2, 3)
	recorder := httptest.NewRecorder()
	app.UserInvite(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	var response struct {
		InvitationID int64 `json:"invitationId"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	var invitationID int64
	var inviterID int
	var status string
	if err := db.QueryRow(`
		SELECT id, inviter_id, status
		FROM group_invitations
		WHERE group_id = 7 AND user_id = 3
	`).Scan(&invitationID, &inviterID, &status); err != nil {
		t.Fatal(err)
	}
	if inviterID != 2 || status != "pending" {
		t.Fatalf("inviter = %d, status = %q; expected inviter 2 and pending", inviterID, status)
	}

	var relatedID int64
	if err := db.QueryRow(`
		SELECT related_id
		FROM notifications
		WHERE user_id = 3 AND actor_id = 2 AND category = 'groups' AND type = 'invitation'
	`).Scan(&relatedID); err != nil {
		t.Fatal(err)
	}
	if relatedID != invitationID {
		t.Fatalf("notification related_id = %d, expected invitation id %d", relatedID, invitationID)
	}
	if response.InvitationID != invitationID {
		t.Fatalf("response invitationId = %d, expected %d", response.InvitationID, invitationID)
	}

	creatorRequest := groupInvitationRequest(http.MethodPost, "/api/groups/7/invitations", "7", 1, 4)
	creatorRecorder := httptest.NewRecorder()
	app.UserInvite(creatorRecorder, creatorRequest)
	if creatorRecorder.Code != http.StatusCreated {
		t.Fatalf("creator invitation status = %d, expected %d", creatorRecorder.Code, http.StatusCreated)
	}
}

func TestUserInviteRejectsInvalidInvitationAttempts(t *testing.T) {
	tests := []struct {
		name      string
		inviterID int
		inviteeID int
		prepare   string
		expected  int
	}{
		{
			name:      "non-member inviter",
			inviterID: 4,
			inviteeID: 3,
			expected:  http.StatusForbidden,
		},
		{
			name:      "existing member invitee",
			inviterID: 2,
			inviteeID: 1,
			expected:  http.StatusConflict,
		},
		{
			name:      "duplicate pending invitation",
			inviterID: 2,
			inviteeID: 3,
			prepare: `
				INSERT INTO group_invitations (group_id, user_id, inviter_id, status)
				VALUES (7, 3, 1, 'pending');
			`,
			expected: http.StatusConflict,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := openInvitationAPITestDB(t, true)
			if test.prepare != "" {
				if _, err := db.Exec(test.prepare); err != nil {
					t.Fatal(err)
				}
			}

			request := groupInvitationRequest(http.MethodPost, "/api/groups/7/invitations", "7", test.inviterID, test.inviteeID)
			recorder := httptest.NewRecorder()
			App{DB: db}.UserInvite(recorder, request)

			if recorder.Code != test.expected {
				t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, test.expected, recorder.Body.String())
			}
		})
	}
}

func TestUserInvitePreservesDeclinedHistoryAndAllowsReinvite(t *testing.T) {
	db := openInvitationAPITestDB(t, true)
	_, err := db.Exec(`
		INSERT INTO group_invitations (id, group_id, user_id, inviter_id, status)
		VALUES (10, 7, 3, 1, 'declined');
	`)
	if err != nil {
		t.Fatal(err)
	}

	request := groupInvitationRequest(http.MethodPost, "/api/groups/7/invitations", "7", 2, 3)
	recorder := httptest.NewRecorder()
	App{DB: db}.UserInvite(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 10 AND status = 'declined'", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE group_id = 7 AND user_id = 3", 2)
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE group_id = 7 AND user_id = 3 AND status = 'pending'", 1)

	var pendingID, relatedID int64
	if err := db.QueryRow(`
		SELECT id
		FROM group_invitations
		WHERE group_id = 7 AND user_id = 3 AND status = 'pending'
	`).Scan(&pendingID); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(`
		SELECT related_id
		FROM notifications
		WHERE user_id = 3 AND type = 'invitation'
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&relatedID); err != nil {
		t.Fatal(err)
	}
	if relatedID != pendingID {
		t.Fatalf("notification related_id = %d, expected pending invitation id %d", relatedID, pendingID)
	}
}

func TestUserInviteRollsBackWhenNotificationCreationFails(t *testing.T) {
	db := openInvitationAPITestDB(t, false)
	request := groupInvitationRequest(http.MethodPost, "/api/groups/7/invitations", "7", 2, 3)
	recorder := httptest.NewRecorder()

	App{DB: db}.UserInvite(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, expected %d", recorder.Code, http.StatusInternalServerError)
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations", 0)
}

func TestGetInviteUsersReturnsOnlyNonMembersWithPendingState(t *testing.T) {
	db := openInvitationAPITestDB(t, false)
	_, err := db.Exec(`
		INSERT INTO group_invitations (id, group_id, user_id, inviter_id, status)
		VALUES
			(10, 7, 3, 1, 'declined'),
			(11, 7, 4, 1, 'accepted'),
			(12, 7, 5, 2, 'pending');
	`)
	if err != nil {
		t.Fatal(err)
	}

	request := inviteUsersRequest("7", 2)
	recorder := httptest.NewRecorder()
	App{DB: db}.GetInviteUsers(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var response struct {
		Users []struct {
			ID           int    `json:"id"`
			IsInvited    bool   `json:"isInvited"`
			InvitationID *int64 `json:"invitationId"`
		} `json:"users"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	users := make(map[int]struct {
		isInvited    bool
		invitationID *int64
	}, len(response.Users))
	for _, user := range response.Users {
		users[user.ID] = struct {
			isInvited    bool
			invitationID *int64
		}{user.IsInvited, user.InvitationID}
	}
	if len(users) != 4 {
		t.Fatalf("invite user count = %d, expected 4", len(users))
	}
	if _, found := users[1]; found {
		t.Fatal("group creator was included in invite users")
	}
	if _, found := users[2]; found {
		t.Fatal("existing group member was included in invite users")
	}
	for _, userID := range []int{3, 4, 6} {
		if user := users[userID]; user.isInvited || user.invitationID != nil {
			t.Fatalf("user %d was marked invited by historical data: %+v", userID, user)
		}
	}
	pendingUser := users[5]
	if !pendingUser.isInvited || pendingUser.invitationID == nil || *pendingUser.invitationID != 12 {
		t.Fatalf("pending user state = %+v, expected invitation 12", pendingUser)
	}

	if _, err = db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (7, 5)"); err != nil {
		t.Fatal(err)
	}
	recorder = httptest.NewRecorder()
	App{DB: db}.GetInviteUsers(recorder, inviteUsersRequest("7", 2))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status after membership = %d, expected %d", recorder.Code, http.StatusOK)
	}
	response.Users = nil
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	for _, user := range response.Users {
		if user.ID == 5 {
			t.Fatal("new group member remained in invite users")
		}
	}
}

func TestGetInviteUsersRejectsNonMember(t *testing.T) {
	db := openInvitationAPITestDB(t, false)
	recorder := httptest.NewRecorder()

	App{DB: db}.GetInviteUsers(recorder, inviteUsersRequest("7", 4))

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, expected %d", recorder.Code, http.StatusForbidden)
	}
}

func TestUndoInvitationDeletesOnlyPendingInvitationAndAllowsReinvite(t *testing.T) {
	db := openInvitationAPITestDB(t, true)
	_, err := db.Exec(`
		INSERT INTO groups (id, creator_id) VALUES (8, 1);
		INSERT INTO group_members (group_id, user_id) VALUES (8, 2);
		INSERT INTO group_invitations (id, group_id, user_id, inviter_id, status)
		VALUES
			(10, 7, 3, 1, 'declined'),
			(11, 7, 3, 1, 'pending'),
			(12, 7, 3, 2, 'accepted');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES
			(20, 3, 1, 'groups', 'invitation', 'old declined invitation', 10),
			(21, 3, 1, 'groups', 'invitation', 'pending invitation', 11),
			(22, 3, 2, 'groups', 'invitation', 'old accepted invitation', 12);
	`)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	App{DB: db}.UndoInvitation(recorder, undoInvitationRequest("8", "11", 2))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("wrong-group status = %d, expected %d", recorder.Code, http.StatusNotFound)
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 11 AND status = 'pending'", 1)

	recorder = httptest.NewRecorder()
	App{DB: db}.UndoInvitation(recorder, undoInvitationRequest("7", "11", 2))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 11", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 10 AND status = 'declined'", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 12 AND status = 'accepted'", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 21", 0)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id IN (20, 22)", 2)

	recorder = httptest.NewRecorder()
	App{DB: db}.UndoInvitation(recorder, undoInvitationRequest("7", "12", 1))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("accepted invitation undo status = %d, expected %d", recorder.Code, http.StatusNotFound)
	}

	recorder = httptest.NewRecorder()
	App{DB: db}.UserInvite(recorder, groupInvitationRequest(http.MethodPost, "/api/groups/7/invitations", "7", 2, 3))
	if recorder.Code != http.StatusCreated {
		t.Fatalf("reinvite status = %d, expected %d; body = %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE group_id = 7 AND user_id = 3 AND status = 'pending'", 1)
}

func TestUndoInvitationRejectsNonMember(t *testing.T) {
	db := openInvitationAPITestDB(t, true)
	_, err := db.Exec(`
		INSERT INTO group_invitations (id, group_id, user_id, inviter_id, status)
		VALUES (11, 7, 3, 1, 'pending');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES (21, 3, 1, 'groups', 'invitation', 'pending invitation', 11);
	`)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	App{DB: db}.UndoInvitation(recorder, undoInvitationRequest("7", "11", 4))

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, expected %d", recorder.Code, http.StatusForbidden)
	}
	assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 11 AND status = 'pending'", 1)
	assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 21", 1)
}

func TestInvitationNotificationActionsUseExactInvitation(t *testing.T) {
	t.Run("accept", func(t *testing.T) {
		db := openInvitationActionTestDB(t)
		request := notificationActionRequestForTest(21, 2, "join")
		recorder := httptest.NewRecorder()

		App{DB: db}.ApplyNotificationAction(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 10 AND status = 'declined'", 1)
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 11 AND status = 'accepted'", 1)
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_members WHERE group_id = 7 AND user_id = 2", 1)
		assertRowCount(t, db, "SELECT COUNT(*) FROM chat_users WHERE chat_id = 30 AND user_id = 2", 1)
		assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 21 AND is_read = 1", 1)
	})

	t.Run("decline", func(t *testing.T) {
		db := openInvitationActionTestDB(t)
		request := notificationActionRequestForTest(21, 2, "decline")
		recorder := httptest.NewRecorder()

		App{DB: db}.ApplyNotificationAction(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
		}
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id IN (10, 11) AND status = 'declined'", 2)
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_members", 0)
		assertRowCount(t, db, "SELECT COUNT(*) FROM notifications WHERE id = 21 AND is_read = 1", 1)
	})

	t.Run("historical notification cannot act on newer invitation", func(t *testing.T) {
		db := openInvitationActionTestDB(t)
		request := notificationActionRequestForTest(20, 2, "join")
		recorder := httptest.NewRecorder()

		App{DB: db}.ApplyNotificationAction(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
		}
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_invitations WHERE id = 11 AND status = 'pending'", 1)
		assertRowCount(t, db, "SELECT COUNT(*) FROM group_members", 0)
	})
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

func openInvitationAPITestDB(t *testing.T, withNotifications bool) *sql.DB {
	t.Helper()

	db := openGroupAPITestDB(t)
	_, err := db.Exec(`
		CREATE TABLE user (
			id INTEGER PRIMARY KEY,
			username TEXT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL
		);
		CREATE TABLE groups (id INTEGER PRIMARY KEY, creator_id INTEGER NOT NULL);
		CREATE TABLE group_members (group_id INTEGER NOT NULL, user_id INTEGER NOT NULL, PRIMARY KEY (group_id, user_id));
		CREATE TABLE group_invitations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			inviter_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE UNIQUE INDEX group_invitations_one_pending
		ON group_invitations (group_id, user_id)
		WHERE status = 'pending';
		INSERT INTO user (id, username, first_name, last_name) VALUES
			(1, 'owner', 'Group', 'Owner'),
			(2, 'member', 'Group', 'Member'),
			(3, 'ahmed', 'Ahmed', 'Ali'),
			(4, 'sara', 'Sara', 'Hassan'),
			(5, 'mohammed', 'Mohammed', 'Saleh'),
			(6, 'moh123', 'Moh', 'Yusuf');
		INSERT INTO groups (id, creator_id) VALUES (7, 1);
		INSERT INTO group_members (group_id, user_id) VALUES (7, 1), (7, 2);
	`)
	if err != nil {
		t.Fatal(err)
	}

	if withNotifications {
		_, err = db.Exec(`
			CREATE TABLE notifications (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				actor_id INTEGER,
				category TEXT NOT NULL,
				type TEXT NOT NULL,
				message TEXT NOT NULL,
				related_id INTEGER,
				is_read INTEGER NOT NULL DEFAULT 0
			);
			CREATE TRIGGER delete_group_invitation_notifications
			AFTER DELETE ON group_invitations
			FOR EACH ROW
			BEGIN
				DELETE FROM notifications
				WHERE category = 'groups'
				  AND type = 'invitation'
				  AND related_id = OLD.id;
			END;
		`)
		if err != nil {
			t.Fatal(err)
		}
	}

	return db
}

func openInvitationActionTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db := openGroupAPITestDB(t)
	_, err := db.Exec(`
		CREATE TABLE user (id INTEGER PRIMARY KEY);
		CREATE TABLE groups (id INTEGER PRIMARY KEY, creator_id INTEGER NOT NULL);
		CREATE TABLE group_members (group_id INTEGER NOT NULL, user_id INTEGER NOT NULL, PRIMARY KEY (group_id, user_id));
		CREATE TABLE chats (id INTEGER PRIMARY KEY, group_id INTEGER NOT NULL);
		CREATE TABLE chat_users (user_id INTEGER NOT NULL, chat_id INTEGER NOT NULL, PRIMARY KEY (user_id, chat_id));
		CREATE TABLE group_join_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT NOT NULL
		);
		CREATE TABLE group_invitations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			inviter_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending'
		);
		CREATE TABLE notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			actor_id INTEGER,
			category TEXT NOT NULL,
			type TEXT NOT NULL,
			message TEXT NOT NULL,
			related_id INTEGER,
			is_read INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		INSERT INTO user (id) VALUES (1), (2);
		INSERT INTO groups (id, creator_id) VALUES (7, 1);
		INSERT INTO chats (id, group_id) VALUES (30, 7);
		INSERT INTO group_invitations (id, group_id, user_id, inviter_id, status)
		VALUES (10, 7, 2, 1, 'declined'), (11, 7, 2, 1, 'pending');
		INSERT INTO notifications (id, user_id, actor_id, category, type, message, related_id)
		VALUES
			(20, 2, 1, 'groups', 'invitation', 'old invitation', 10),
			(21, 2, 1, 'groups', 'invitation', 'new invitation', 11);
	`)
	if err != nil {
		t.Fatal(err)
	}

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

func groupInvitationRequest(method string, path string, groupID string, inviterID int, inviteeID int) *http.Request {
	body := strings.NewReader(`{"userId":` + strconv.Itoa(inviteeID) + `}`)
	request := httptest.NewRequest(method, path, body)
	request.SetPathValue("id", groupID)
	ctx := context.WithValue(request.Context(), "userID", inviterID)
	return request.WithContext(ctx)
}

func inviteUsersRequest(groupID string, userID int) *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/api/groups/"+groupID+"/invite-users", nil)
	request.SetPathValue("id", groupID)
	ctx := context.WithValue(request.Context(), "userID", userID)
	return request.WithContext(ctx)
}

func undoInvitationRequest(groupID string, invitationID string, userID int) *http.Request {
	request := httptest.NewRequest(http.MethodDelete, "/api/groups/"+groupID+"/invitations/"+invitationID, nil)
	request.SetPathValue("id", groupID)
	request.SetPathValue("invitationID", invitationID)
	ctx := context.WithValue(request.Context(), "userID", userID)
	return request.WithContext(ctx)
}

func notificationActionRequestForTest(notificationID int64, userID int, action string) *http.Request {
	body := strings.NewReader(`{"action":"` + action + `"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/notifications/action", body)
	request.SetPathValue("notificationID", strconv.FormatInt(notificationID, 10))
	ctx := context.WithValue(request.Context(), "userID", userID)
	return request.WithContext(ctx)
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

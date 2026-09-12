package api

import (
	"encoding/json"
	"net/http"
	"social/database/events"
	"social/internal/helpers"
	"strconv"
	"strings"
	"time"
)

func (app App) EventRSVP(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	id, err := strconv.ParseInt(r.PathValue("eventID"), 10, 64)
	var input struct {
		Response string `json:"response"`
	}
	if err != nil || json.NewDecoder(http.MaxBytesReader(w, r.Body, 2048)).Decode(&input) != nil {
		helpers.WriteJson(w, 400, map[string]any{"message": "invalid response"})
		return
	}
	if err := events.SetRSVP(app.DB, userID, id, input.Response); err != nil {
		helpers.WriteJson(w, 400, map[string]any{"message": err.Error()})
		return
	}
	helpers.WriteJson(w, 200, map[string]any{"status": true})
}

func (app App) groupMember(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	userID, ok := r.Context().Value("userID").(int)
	groupID, err := strconv.Atoi(r.PathValue("id"))
	if !ok || err != nil || groupID <= 0 {
		helpers.WriteJson(w, 400, map[string]any{"message": "invalid group"})
		return 0, 0, false
	}
	var count int
	err = app.DB.QueryRow(`SELECT COUNT(*) FROM group_members WHERE group_id=? AND user_id=?`, groupID, userID).Scan(&count)
	if err != nil || count != 1 {
		helpers.WriteJson(w, 403, map[string]any{"message": "group membership required"})
		return 0, 0, false
	}
	return userID, groupID, true
}

func (app App) InviteGroupMember(w http.ResponseWriter, r *http.Request) {
	userID, groupID, ok := app.groupMember(w, r)
	if !ok {
		return
	}
	var input struct {
		UserID int `json:"userId"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 2048)).Decode(&input) != nil || input.UserID <= 0 || input.UserID == userID {
		helpers.WriteJson(w, 400, map[string]any{"message": "choose another user"})
		return
	}
	tx, err := app.DB.Begin()
	if err != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not invite user"})
		return
	}
	defer tx.Rollback()
	var available int
	err = tx.QueryRow(`SELECT COUNT(*) FROM user WHERE id=? AND NOT EXISTS (SELECT 1 FROM group_members WHERE group_id=? AND user_id=?)`, input.UserID, groupID, input.UserID).Scan(&available)
	if err != nil || available != 1 {
		helpers.WriteJson(w, 400, map[string]any{"message": "user not found or already a member"})
		return
	}
	result, err := tx.Exec(`INSERT INTO group_invitations (group_id,user_id,inviter_id) VALUES (?,?,?) ON CONFLICT(group_id,user_id) DO NOTHING`, groupID, input.UserID, userID)
	if err != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not save invitation"})
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		helpers.WriteJson(w, 409, map[string]any{"message": "this user has already been invited"})
		return
	}
	_, err = tx.Exec(`INSERT INTO notifications (user_id,actor_id,category,type,message,related_id) SELECT ?,?,'groups','invitation','You are invited to ' || title,id FROM groups WHERE id=?`, input.UserID, userID, groupID)
	if err != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not save notification"})
		return
	}
	if tx.Commit() != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not finish invitation"})
		return
	}
	helpers.WriteJson(w, 201, map[string]any{"status": true})
}

func (app App) GroupEvents(w http.ResponseWriter, r *http.Request) {
	userID, groupID, ok := app.groupMember(w, r)
	if !ok {
		return
	}
	if r.Method == http.MethodGet {
		rows, err := app.DB.Query(`SELECT e.id,e.title,e.content,COALESCE(e.starts_at,''),COALESCE(v.response,'') FROM events e LEFT JOIN event_rsvps v ON v.event_id=e.id AND v.user_id=? WHERE e.group_id=? ORDER BY e.starts_at,e.id`, userID, groupID)
		if err != nil {
			helpers.WriteJson(w, 500, map[string]any{"message": "could not load events"})
			return
		}
		defer rows.Close()
		items := []map[string]any{}
		for rows.Next() {
			var id int
			var title, description, startsAt, response string
			if rows.Scan(&id, &title, &description, &startsAt, &response) != nil {
				helpers.WriteJson(w, 500, map[string]any{"message": "could not read events"})
				return
			}
			items = append(items, map[string]any{"id": id, "title": title, "description": description, "startsAt": startsAt, "response": response})
		}
		if rows.Err() != nil {
			helpers.WriteJson(w, 500, map[string]any{"message": "could not read events"})
			return
		}
		helpers.WriteJson(w, 200, map[string]any{"events": items})
		return
	}
	var input struct {
		Title       string
		Description string
		StartsAt    string
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&input) != nil {
		helpers.WriteJson(w, 400, map[string]any{"message": "invalid event"})
		return
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	startsAt, err := time.Parse(time.RFC3339, input.StartsAt)
	if err != nil || !startsAt.After(time.Now()) || len([]rune(input.Title)) == 0 || len([]rune(input.Title)) > 50 || len([]rune(input.Description)) == 0 || len([]rune(input.Description)) > 500 {
		helpers.WriteJson(w, 400, map[string]any{"message": "provide a title, description and future event time"})
		return
	}
	tx, err := app.DB.Begin()
	if err != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not create event"})
		return
	}
	defer tx.Rollback()
	result, err := tx.Exec(`INSERT INTO events (title,content,group_id,creator_id,starts_at) VALUES (?,?,?,?,?)`, input.Title, input.Description, groupID, userID, startsAt.UTC().Format(time.RFC3339))
	if err != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not save event"})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not save event"})
		return
	}
	_, err = tx.Exec(`INSERT INTO notifications (user_id,actor_id,category,type,message,related_id) SELECT user_id,?,'events','event_created',?,? FROM group_members WHERE group_id=? AND user_id<>?`, userID, "New event: "+input.Title, id, groupID, userID)
	if err != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not notify members"})
		return
	}
	if tx.Commit() != nil {
		helpers.WriteJson(w, 500, map[string]any{"message": "could not finish event"})
		return
	}
	helpers.WriteJson(w, 201, map[string]any{"status": true, "id": id})
}

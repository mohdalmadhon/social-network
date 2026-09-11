package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"social/database/notifications"
	"social/internal/helpers"
	"social/internal/models"
	"strconv"

	"golang.org/x/net/websocket"
)

func (app *App) HandleWS(ws *websocket.Conn) {
	userID, ok := ws.Request().Context().Value("userID").(int)

	if !ok {
		ws.Close()
		return
	}

	app.Mu.Lock()
	app.Conns[userID] = ws
	app.Mu.Unlock()

	defer func() {
		app.Mu.Lock()
		delete(app.Conns, userID)
		app.Mu.Unlock()
		ws.Close()
	}()

	app.readLoop(ws)
}

func (app *App) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 4096)

	for {
		n, err := ws.Read(buff)

		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println(err)
			break
		}

		var payload models.WSPayload

		if err := json.Unmarshal(buff[:n], &payload); err != nil {
			log.Println(err)
			continue
		}

		switch payload.Type {
		case "message":
			app.handleMessage(payload.Data)

		case "notification":
			app.handleNotification(payload.Data)

		default:
			log.Println("unknown websocket type:", payload.Type)
		}
	}
}

func (app *App) SendToUser(userID int, msg []byte) {
	app.Mu.Lock()
	ws, ok := app.Conns[userID]
	app.Mu.Unlock()

	if !ok {
		return
	}

	if _, err := ws.Write(msg); err != nil {
		log.Println("websocket failed:", err)
	}
}

func (app *App) handleMessage(data json.RawMessage) {

}

func (app *App) handleNotification(data json.RawMessage) {
	var notification models.NewNotification

	if err := json.Unmarshal(data, &notification); err != nil {
		log.Println(err)
		return
	}

	if err := notifications.InsertNotification(app.DB, notification); err != nil {
		log.Println("failed to insert notification:", err)
		return
	}

	msg, err := json.Marshal(models.WSPayload{
		Type: "notification",
		Data: data,
	})

	if err != nil {
		log.Println(err)
		return
	}

	app.SendToUser(notification.UserID, msg)
}

func (app *App) GetNotification(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	offset := 0

	if value := r.URL.Query().Get("offset"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	rows, err := app.DB.Query(`
		SELECT
			n.id,
			n.message,
			n.created_at,
			nt.notifications_id,
			nt.post_id_tag,
			nt.message_user_id,
			nt.comment_reply_user_id,
			nt.follow_request_user_id,
			nt.follow_request_accept_user_id,
			nt.follow_user_id,
			nt.post_like_user_id,
			nt.post_dislike_user_id,
			nt.comment_like_user_id,
			nt.comment_mention_user_id,
			nt.post_mention_user_id,
			nt.group_invite_user_id,
			nt.group_join_user_id,
			nt.group_accept_user_id,
			nt.event_invite_user_id,
			nt.event_response_user_id
		FROM notifications n
		LEFT JOIN notifications_types nt
			ON nt.notifications_id = n.id
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC
		LIMIT 20 OFFSET ?
	`, userID, offset)

	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get notifications",
		})
		return
	}
	defer rows.Close()

	notifications := []map[string]any{}

	for rows.Next() {
		var (
			id                        int
			message                   string
			createdAt                 string
			notificationsID           int
			postIDTag                 sql.NullInt64
			messageUserID             sql.NullInt64
			commentReplyUserID        sql.NullInt64
			followRequestUserID       sql.NullInt64
			followRequestAcceptUserID sql.NullInt64
			followUserID              sql.NullInt64
			postLikeUserID            sql.NullInt64
			postDislikeUserID         sql.NullInt64
			commentLikeUserID         sql.NullInt64
			commentMentionUserID      sql.NullInt64
			postMentionUserID         sql.NullInt64
			groupInviteUserID         sql.NullInt64
			groupJoinUserID           sql.NullInt64
			groupAcceptUserID         sql.NullInt64
			eventInviteUserID         sql.NullInt64
			eventResponseUserID       sql.NullInt64
		)

		err := rows.Scan(
			&id,
			&message,
			&createdAt,
			&notificationsID,
			&postIDTag,
			&messageUserID,
			&commentReplyUserID,
			&followRequestUserID,
			&followRequestAcceptUserID,
			&followUserID,
			&postLikeUserID,
			&postDislikeUserID,
			&commentLikeUserID,
			&commentMentionUserID,
			&postMentionUserID,
			&groupInviteUserID,
			&groupJoinUserID,
			&groupAcceptUserID,
			&eventInviteUserID,
			&eventResponseUserID,
		)

		if err != nil {
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not read notifications",
			})
			return
		}

		notification := map[string]any{
			"id":         id,
			"message":    message,
			"created_at": createdAt,
		}

		if postIDTag.Valid {
			notification["post_id"] = postIDTag.Int64
		}

		if messageUserID.Valid {
			notification["message_user_id"] = messageUserID.Int64
		}

		if commentReplyUserID.Valid {
			notification["comment_reply_user_id"] = commentReplyUserID.Int64
		}

		if followRequestUserID.Valid {
			notification["follow_request_user_id"] = followRequestUserID.Int64
		}

		if followRequestAcceptUserID.Valid {
			notification["follow_request_accept_user_id"] = followRequestAcceptUserID.Int64
		}

		if followUserID.Valid {
			notification["follow_user_id"] = followUserID.Int64
		}

		if postLikeUserID.Valid {
			notification["post_like_user_id"] = postLikeUserID.Int64
		}

		if postDislikeUserID.Valid {
			notification["post_dislike_user_id"] = postDislikeUserID.Int64
		}

		if commentLikeUserID.Valid {
			notification["comment_like_user_id"] = commentLikeUserID.Int64
		}

		if commentMentionUserID.Valid {
			notification["comment_mention_user_id"] = commentMentionUserID.Int64
		}

		if postMentionUserID.Valid {
			notification["post_mention_user_id"] = postMentionUserID.Int64
		}

		if groupInviteUserID.Valid {
			notification["group_invite_user_id"] = groupInviteUserID.Int64
		}

		if groupJoinUserID.Valid {
			notification["group_join_user_id"] = groupJoinUserID.Int64
		}

		if groupAcceptUserID.Valid {
			notification["group_accept_user_id"] = groupAcceptUserID.Int64
		}

		if eventInviteUserID.Valid {
			notification["event_invite_user_id"] = eventInviteUserID.Int64
		}

		if eventResponseUserID.Valid {
			notification["event_response_user_id"] = eventResponseUserID.Int64
		}

		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not read notifications",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":        true,
		"notifications": notifications,
		"offset":        offset,
		"limit":         20,
		"hasMore":       len(notifications) == 20,
	})
}

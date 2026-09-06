package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"social/database/events"
	"social/database/groups"
	"social/database/notifications"
	profiles "social/database/profile"
	"social/internal/models"
	"strconv"
)

func (app App) Notifications(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
		return
	}

	category := r.URL.Query().Get("category")
	result, err := notifications.List(app.DB, userID, category)
	if errors.Is(err, notifications.ErrInvalidCategory) {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not load notifications",
		})
		return
	}

	unreadCount, err := notifications.UnreadCount(app.DB, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not load notification count",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":        true,
		"notifications": result,
		"unreadCount":   unreadCount,
	})
}

func (app App) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", http.MethodPatch)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
		return
	}

	notificationID, err := strconv.ParseInt(r.PathValue("notificationID"), 10, 64)
	if err != nil || notificationID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid notification id",
		})
		return
	}

	if err := notifications.MarkRead(app.DB, userID, notificationID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not mark notification as read",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": true})
}

func (app App) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", http.MethodPatch)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
		return
	}

	if err := notifications.MarkAllRead(app.DB, userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not mark notifications as read",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": true})
}

type notificationActionRequest struct {
	Action string `json:"action"`
}

func (app App) ApplyNotificationAction(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", http.MethodPatch)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
		return
	}

	notificationID, err := strconv.ParseInt(r.PathValue("notificationID"), 10, 64)
	if err != nil || notificationID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid notification id",
		})
		return
	}

	var request notificationActionRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid notification action",
		})
		return
	}

	notification, err := notifications.GetByID(app.DB, userID, notificationID)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "notification not found",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not load notification",
		})
		return
	}

	if err := app.applyNotificationAction(userID, notification, request.Action); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, groups.ErrGroupNotFound) || errors.Is(err, events.ErrEventNotFound) {
			status = http.StatusNotFound
		}
		writeJSON(w, status, map[string]any{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	if err := notifications.MarkRead(app.DB, userID, notificationID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "action completed but notification could not be marked read",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": true,
		"action": request.Action,
	})
}

func (app App) applyNotificationAction(userID int, notification models.Notification, action string) error {
	switch notification.Category {
	case "requests":
		if notification.Type != "follow_request" || notification.ActorID == nil {
			return errors.New("this request cannot be acted on")
		}
		switch action {
		case "accept":
			return profiles.SendFollowRequest(app.DB, userID, *notification.ActorID, 1)
		case "decline":
			return profiles.SendFollowRequest(app.DB, userID, *notification.ActorID, -1)
		default:
			return errors.New("request action must be accept or decline")
		}

	case "groups":
		if action != "join" {
			return errors.New("group action must be join")
		}
		if notification.RelatedID == nil {
			return errors.New("group notification is missing its group")
		}
		return groups.JoinGroup(app.DB, userID, *notification.RelatedID)

	case "events":
		if notification.RelatedID == nil {
			return errors.New("event notification is missing its event")
		}
		switch action {
		case "rsvp":
			return events.SetRSVP(app.DB, userID, *notification.RelatedID, "going")
		case "decline":
			return events.SetRSVP(app.DB, userID, *notification.RelatedID, "declined")
		default:
			return errors.New("event action must be rsvp or decline")
		}

	default:
		return errors.New("unsupported notification category")
	}
}

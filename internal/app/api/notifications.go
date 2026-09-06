package api

import (
	"errors"
	"net/http"
	"social/database/notifications"
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

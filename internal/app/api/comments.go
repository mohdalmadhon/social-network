package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"social/database/comments"
	"social/internal/models"
	"strconv"
	"strings"
)

func (app App) Comments(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	postID, err := strconv.ParseInt(r.PathValue("postID"), 10, 64)
	if err != nil || postID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid post id",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.listComments(w, userID, postID)
	case http.MethodPost:
		app.createComment(w, r, userID, postID)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
	}
}

func (app App) listComments(w http.ResponseWriter, userID int, postID int64) {
	result, err := comments.ListComments(app.DB, userID, postID)
	if errors.Is(err, comments.ErrPostNotVisible) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "post is not available",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not load comments",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":   true,
		"comments": result,
	})
}

func (app App) createComment(w http.ResponseWriter, r *http.Request, userID int, postID int64) {
	var request models.CreateCommentRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid comment body",
		})
		return
	}

	request.Content = strings.TrimSpace(request.Content)
	if request.Content == "" || len([]rune(request.Content)) > 200 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "comment must contain 1 to 200 characters",
		})
		return
	}

	comment, err := comments.CreateComment(app.DB, userID, postID, request.Content)
	if errors.Is(err, comments.ErrPostNotVisible) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "post is not available",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not create comment",
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status":  true,
		"comment": comment,
	})
}

package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"social/database/posts"
	"social/internal/app/tokens"
	"social/internal/models"
	"strings"
)

const maxPostBodySize = 1 << 20

func (app *App) Posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.listPosts(w, r)
	case http.MethodPost:
		app.createPost(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
	}
}

func (app App) createPost(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxPostBodySize)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request models.CreatePostRequest
	if err = decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid JSON body",
		})
		return
	}
	if err = decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "request body must contain one JSON object",
		})
		return
	}

	request.Content = strings.TrimSpace(request.Content)
	if request.Content == "" || len([]rune(request.Content)) > 500 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "post content must contain 1 to 500 characters",
		})
		return
	}
	if !posts.IsPostPrivacy(request.Privacy) {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "privacy must be public, followers, or selected",
		})
		return
	}
	if err = posts.ValidateSelectedIDs(request.Privacy, request.SelectedFollowerIDs); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	post, err := posts.CreatePost(app.DB, userID, request)
	if errors.Is(err, posts.ErrSelectedFollowersRequired) || errors.Is(err, posts.ErrInvalidPostViewer) {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not create post",
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status": true,
		"post":   post,
	})
}

func (app App) listPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	posts, err := posts.ListFeedPosts(app.DB, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not load feed",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": true,
		"posts":  posts,
	})
}

func authenticatedUserID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		return 0, errors.New("missing authentication cookie")
	}

	payload, err := tokens.VerifyToken(cookie.Value)
	if err != nil {
		return 0, err
	}
	if payload.UserID <= 0 {
		return 0, errors.New("invalid authenticated user")
	}

	return payload.UserID, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

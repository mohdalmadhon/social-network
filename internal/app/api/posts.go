package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"social/database/posts"
	"social/internal/app/tokens"
	"social/internal/helpers"
	"social/internal/models"
	"strings"
)

const maxPostBodySize = 6 << 20
const maxPostImageSize = 5 * 1024 * 1024

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
	if err = r.ParseMultipartForm(maxPostBodySize); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "request body too large",
		})
		return
	}

	request := models.CreatePostRequest{
		Content: r.FormValue("content"),
		Privacy: r.FormValue("privacy"),
	}
	selectedIDs := r.FormValue("selectedFollowerIds")
	if selectedIDs != "" {
		if err := json.Unmarshal([]byte(selectedIDs), &request.SelectedFollowerIDs); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid selected followers",
			})
			return
		}
	}

	imageFile, imageHeader, fileErr := r.FormFile("image")

	if fileErr != nil && fileErr != http.ErrMissingFile {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid image upload",
		})
		return
	}

	if fileErr == nil {
		defer imageFile.Close()

		if imageHeader.Size > maxPostImageSize {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "image must be smaller than 5 MB",
			})
			return
		}

		fileBytes := make([]byte, 512)
		bytesRead, readErr := imageFile.Read(fileBytes)

		if readErr != nil && readErr != io.EOF {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "could not read image",
			})
			return
		}

		contentType := http.DetectContentType(fileBytes[:bytesRead])

		if contentType != "image/jpeg" &&
			contentType != "image/png" &&
			contentType != "image/gif" {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "only JPEG, PNG, and GIF images are allowed",
			})
			return
		}

		if _, seekErr := imageFile.Seek(0, io.SeekStart); seekErr != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "could not reset image",
			})
			return
		}
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

	if fileErr == nil {
		imagePath, saveErr := helpers.SaveUploads(imageFile, imageHeader, "post")
		if saveErr != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not save image",
			})
			return
		}

		request.ImagePath = imagePath
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

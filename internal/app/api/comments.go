package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"social/database/comments"
	"social/database/posts"
	"social/internal/helpers"
	"social/internal/models"
	"strconv"
	"strings"
)

var (
	ErrPostNotVisible  = errors.New("post not visible")
	ErrCommentNotFound = errors.New("comment not found")
)

const maxCommentBodySize = 4 << 10
const maxCommentMultipartSize = 6 << 20
const maxCommentImageSize = 5 * 1024 * 1024

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
		app.listComments(w, r, userID, postID)
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

func (app App) listComments(w http.ResponseWriter, r *http.Request, userID int, postID int64) {
	page, err := parsePage(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "limit must be between 1 and 50 and offset cannot be negative",
		})
		return
	}

	result, err := comments.ListComments(app.DB, userID, postID, page.Limit+1, page.Offset)
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
	result, hasMore := trimPage(result, page, true)
	for i := range result {
		result[i].Own = result[i].UserID == userID
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":     true,
		"comments":   result,
		"hasMore":    hasMore,
		"nextOffset": page.Offset + len(result),
	})
}

func (app App) createComment(w http.ResponseWriter, r *http.Request, userID int, postID int64) {
	var request models.CreateCommentRequest
	var imagePath string

	if strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, maxCommentMultipartSize)
		if err := r.ParseMultipartForm(maxCommentMultipartSize); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "request body too large",
			})
			return
		}

		request.Content = r.FormValue("content")

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

			if imageHeader.Size > maxCommentImageSize {
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

			savedPath, saveErr := helpers.SaveUploads(imageFile, imageHeader, "comment")
			if saveErr != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]any{
					"status":  false,
					"message": "could not save image",
				})
				return
			}

			imagePath = savedPath
		}
	} else {
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxCommentBodySize))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid comment body",
			})
			return
		}
	}

	request.Content = strings.TrimSpace(request.Content)
	if request.Content == "" && imagePath == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "comment must contain text, an image, or both",
		})
		return
	}
	if len([]rune(request.Content)) > 200 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "comment text must contain 1 to 200 characters",
		})
		return
	}

	comment, err := comments.CreateComment(app.DB, userID, postID, request.Content, imagePath)
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

	comment.Own = true

	writeJSON(w, http.StatusCreated, map[string]any{
		"status":  true,
		"comment": comment,
	})
}

func (app App) DeleteComment(w http.ResponseWriter, r *http.Request) {
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

	commentID, err := strconv.ParseInt(r.PathValue("commentID"), 10, 64)
	if err != nil || commentID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid comment id",
		})
		return
	}

	err = posts.DeleteComment(app.DB, userID, postID, commentID)
	if errors.Is(err, posts.ErrCommentNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "comment not found",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not delete comment",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "comment deleted",
	})
}

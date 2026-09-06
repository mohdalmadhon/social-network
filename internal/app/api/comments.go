package api

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"social/database/comments"
	"social/internal/helpers"
	"social/internal/models"
	"strconv"
	"strings"
)

const maxCommentBodySize = 6 << 20
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
	var imageFile multipart.File
	var imageHeader *multipart.FileHeader

	if strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, maxCommentBodySize)
		if err := r.ParseMultipartForm(maxCommentBodySize); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "comment body is too large or invalid",
			})
			return
		}

		request.Content = r.FormValue("content")
		var fileErr error
		imageFile, imageHeader, fileErr = r.FormFile("image")
		if fileErr != nil && fileErr != http.ErrMissingFile {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"status":  false,
				"message": "invalid comment image upload",
			})
			return
		}

		if imageFile != nil {
			if imageHeader.Size > maxCommentImageSize {
				imageFile.Close()
				writeJSON(w, http.StatusBadRequest, map[string]any{
					"status":  false,
					"message": "comment image must be smaller than 5 MB",
				})
				return
			}

			fileBytes := make([]byte, 512)
			bytesRead, readErr := imageFile.Read(fileBytes)
			if readErr != nil && readErr != io.EOF {
				imageFile.Close()
				writeJSON(w, http.StatusBadRequest, map[string]any{
					"status":  false,
					"message": "could not read comment image",
				})
				return
			}

			contentType := http.DetectContentType(fileBytes[:bytesRead])
			if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
				imageFile.Close()
				writeJSON(w, http.StatusBadRequest, map[string]any{
					"status":  false,
					"message": "only JPEG, PNG, and GIF comment images are allowed",
				})
				return
			}

			if _, seekErr := imageFile.Seek(0, io.SeekStart); seekErr != nil {
				imageFile.Close()
				writeJSON(w, http.StatusBadRequest, map[string]any{
					"status":  false,
					"message": "could not reset comment image",
				})
				return
			}
		}
	} else {
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10))
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
	if (request.Content == "" && imageFile == nil) || len([]rune(request.Content)) > 200 {
		if imageFile != nil {
			imageFile.Close()
		}
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "comment needs text or an image, with text limited to 200 characters",
		})
		return
	}

	imagePath := ""
	var err error
	if imageFile != nil {
		defer imageFile.Close()
		imagePath, err = helpers.SaveUploads(imageFile, imageHeader, "comment")
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not save comment image",
			})
			return
		}
	}

	comment, err := comments.CreateCommentWithImage(app.DB, userID, postID, request.Content, imagePath)
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

package api

import (
	"errors"
	"net/http"
	"social/database/posts"
	"strconv"
)

func (app App) LikePost(w http.ResponseWriter, r *http.Request) {
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

	var result posts.ReactionResult
	if r.Method == http.MethodPut {
		result, err = posts.LikePost(app.DB, userID, postID)
	} else if r.Method == http.MethodDelete {
		result, err = posts.UnlikePost(app.DB, userID, postID)
	} else {
		w.Header().Set("Allow", "PUT, DELETE")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"status":  false,
			"message": "method not allowed",
		})
		return
	}

	if errors.Is(err, posts.ErrPostNotVisible) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "post is not available",
		})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not update post reaction",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":    true,
		"liked":     result.Liked,
		"likeCount": result.LikeCount,
	})
}

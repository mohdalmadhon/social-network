package api

import (
	"encoding/json"
	"net/http"
	"social/database/posts"
	"social/internal/helpers"
	"social/internal/models"
	"social/internal/validation"
)

func (app *App) AddComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	var comment models.Comment
	comment.User.ID = userID
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid comment",
		})
		return
	}

	if err := validation.ValidateComment(comment); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	if err := posts.InsertComment(app.DB, comment); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to insert comment",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "comment inserted!",
	})
}
func (app *App) DeleteComment(w http.ResponseWriter, r *http.Request) {
}

func (app *App) GetComments(w http.ResponseWriter, r *http.Request) {
}

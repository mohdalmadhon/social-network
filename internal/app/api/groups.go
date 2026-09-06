package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (app App) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid request body",
		})
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "group title is required",
		})
		return
	}

	if input.Description == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "group description is required",
		})
		return
	}

	tx, err := app.DB.Begin()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to create group",
		})
		return
	}

	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO groups (creator_id, title, description)
		VALUES (?, ?, ?)
	`, userID, input.Title, input.Description)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to create group",
		})
		return
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to get created group",
		})
		return
	}

	// The creator should automatically be a member of the group.
	_, err = tx.Exec(`
		INSERT INTO group_members (group_id, user_id)
		VALUES (?, ?)
	`, groupID, userID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to add creator to group",
		})
		return
	}

	err = tx.Commit()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to create group",
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status":  true,
		"message": "group created successfully",
		"group": map[string]any{
			"id":          groupID,
			"title":       input.Title,
			"description": input.Description,
			"creator_id":  userID,
		},
	})
}

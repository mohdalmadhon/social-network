package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

func (app App) GetGroups(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	rows, err := app.DB.Query(`
	SELECT
		g.id,
		g.title,
		g.description,
		g.creator_id,
		COUNT(gm.user_id) AS member_count,

		EXISTS (
			SELECT 1
			FROM group_members
			WHERE group_id = g.id
			  AND user_id = ?
		) AS is_member,

		EXISTS (
			SELECT 1
			FROM group_join_requests
			WHERE group_id = g.id
			  AND user_id = ?
			  AND status = 'pending'
		) AS is_requested

	FROM groups g

	LEFT JOIN group_members gm
		ON gm.group_id = g.id

	GROUP BY
		g.id,
		g.title,
		g.description,
		g.creator_id

	ORDER BY g.id DESC
	`, userID, userID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to get groups",
		})
		return
	}
	defer rows.Close()

	type Group struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatorID   int    `json:"creatorId"`
		MemberCount int    `json:"memberCount"`
		IsMember    bool   `json:"isMember"`
		IsRequested bool   `json:"isRequested"`
	}

	groups := []Group{}

	for rows.Next() {
		var group Group

		err := rows.Scan(
			&group.ID,
			&group.Title,
			&group.Description,
			&group.CreatorID,
			&group.MemberCount,
			&group.IsMember,
			&group.IsRequested,
		)

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "failed to read groups",
			})
			return
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to read groups",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": true,
		"groups": groups,
	})
}

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
			"memberCount": 1,
			"isMember":    true,
			"isRequested": false,
		},
	})
}

func (app App) GetGroup(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	groupID := r.PathValue("id")
	if groupID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "group id is required",
		})
		return
	}

	type Group struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CreatorID   int    `json:"creatorId"`
		MemberCount int    `json:"memberCount"`
		IsMember    bool   `json:"isMember"`
		IsRequested bool   `json:"isRequested"`
		IsCreator   bool   `json:"isCreator"`
	}

	var group Group

	err = app.DB.QueryRow(`
		SELECT
			g.id,
			g.title,
			g.description,
			g.creator_id,

			COUNT(gm.user_id) AS member_count,

			EXISTS (
				SELECT 1
				FROM group_members
				WHERE group_id = g.id
				  AND user_id = ?
			) AS is_member,

			EXISTS (
				SELECT 1
				FROM group_join_requests
				WHERE group_id = g.id
				  AND user_id = ?
				  AND status = 'pending'
			) AS is_requested,

			g.creator_id = ? AS is_creator

		FROM groups g

		LEFT JOIN group_members gm
			ON gm.group_id = g.id

		WHERE g.id = ?

		GROUP BY
			g.id,
			g.title,
			g.description,
			g.creator_id
	`, userID, userID, userID, groupID).Scan(
		&group.ID,
		&group.Title,
		&group.Description,
		&group.CreatorID,
		&group.MemberCount,
		&group.IsMember,
		&group.IsRequested,
		&group.IsCreator,
	)

	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"status":  false,
			"message": "group not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to get group",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": true,
		"group":  group,
	})
}

func (app App) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	groupID := r.PathValue("id")
	if groupID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "group id is required",
		})
		return
	}

	result, err := app.DB.Exec(`
		DELETE FROM groups
		WHERE id = ?
		  AND creator_id = ?
	`, groupID, userID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to delete group",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to delete group",
		})
		return
	}

	if rowsAffected == 0 {
		writeJSON(w, http.StatusForbidden, map[string]any{
			"status":  false,
			"message": "group not found or you are not the creator",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "group deleted successfully",
	})
}

func (app App) JoinRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	groupID := r.PathValue("id")

	var isMember bool
	err = app.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM group_members
			WHERE group_id = ?
			  AND user_id = ?
		)
	`, groupID, userID).Scan(&isMember)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to check membership",
		})
		return
	}

	if isMember {
		writeJSON(w, http.StatusConflict, map[string]any{
			"status":  false,
			"message": "you are already a member",
		})
		return
	}

	_, err = app.DB.Exec(`
		INSERT INTO group_join_requests (group_id, user_id, status)
		VALUES (?, ?, 'pending')
		ON CONFLICT(group_id, user_id)
		DO UPDATE SET status = 'pending'
	`, groupID, userID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to send join request",
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status":  true,
		"message": "join request sent",
	})
}

func (app App) UndoJoinRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "authentication required",
		})
		return
	}

	groupID := r.PathValue("id")

	_, err = app.DB.Exec(`
		DELETE FROM group_join_requests
		WHERE group_id = ?
		  AND user_id = ?
		  AND status = 'pending'
	`, groupID, userID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to undo join request",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": true,
	})
}

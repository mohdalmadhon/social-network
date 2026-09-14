package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social/database/chats"
	"social/internal/helpers"
	"strconv"
)

func (app *App) GetGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	privatChatQuery := r.URL.Query().Get("private")
	isPrivate, err := strconv.Atoi(privatChatQuery)
	if err != nil || (isPrivate != 0 && isPrivate != 1) {
		helpers.WriteJson(w, http.StatusBadGateway, map[string]any{
			"status":  false,
			"message": "what you want to get excatly? specify please.........  bro....",
		})
		return
	}

	offset := 0
	offsetQuery := r.URL.Query().Get("offset")
	if offsetQuery != "" && offsetQuery != "null" && offsetQuery != "undifiend" {
		var err error
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil {
			helpers.WriteJson(w, http.StatusBadGateway, map[string]any{
				"status":  false,
				"message": "invalid offset",
			})
			return
		}
	}

	chats, err := chats.GetPrivateChatsList(app.DB, userID, offset)
	if err != nil && err != sql.ErrNoRows {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "failed to get users list",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status": true,
		"data":   chats,
	})
}

func (app *App) SearchPrivateChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}
	privatChatQuery := r.URL.Query().Get("private")
	isPrivate, err := strconv.Atoi(privatChatQuery)
	if err != nil || (isPrivate != 0 && isPrivate != 1) {
		helpers.WriteJson(w, http.StatusBadGateway, map[string]any{
			"status":  false,
			"message": "what you want to get excatly? specify please.........  bro....",
		})
		return
	}

	offset := 0
	offsetQuery := r.URL.Query().Get("offset")
	if offsetQuery != "" && offsetQuery != "null" && offsetQuery != "undifiend" {
		var err error
		offset, err = strconv.Atoi(offsetQuery)
		if err != nil {
			helpers.WriteJson(w, http.StatusBadGateway, map[string]any{
				"status":  false,
				"message": "invalid offset",
			})
			return
		}
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status": true,
			"data":   nil,
		})
		return
	}

	if isPrivate == 1 {
		chats, err := chats.SearchChatUsers(app.DB, userID, offset, search)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "failed to get chats",
			})
			return
		}

		helpers.WriteJson(w, http.StatusOK, map[string]any{
			"status": true,
			"data":   chats,
		})
	}
}

func (app *App) AddMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	type request struct {
		Offset  int    `json:"offset"`
		GroupID int    `json:"groupID"`
		UserID  int    `json:"userID"`
		Content string `json:"content"`
	}

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid data",
		})
		return
	}

	if req.Content == "" {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "message cannot be empty",
		})
		return
	}

	if req.GroupID != 1 {
		// group message
	}
	log.Println(req.UserID)
	groupID := req.GroupID

	if groupID <= 0 {
		existingGroupID, err := chats.HasPrivateChat(app.DB, userID, req.UserID)
		if err != nil {
			helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
				"status":  false,
				"message": "could not verify chat",
			})
			return
		}

		if existingGroupID != -1 {
			groupID = existingGroupID
		} else {

			groupID, err = chats.MakePrivateChat(app.DB, userID, req.UserID)
			if err != nil {
				helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
					"status":  false,
					"message": "could not make chat",
				})
				return
			}
		}
	}

	err := chats.AddMessages(app.DB, req.Content, userID, groupID)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not send chat",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status":  true,
		"message": "message sent",
		"groupID": groupID,
	})
}

func (app *App) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		helpers.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"status":  false,
			"message": "could not authorize user",
		})
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("groupID"))
	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid groupID",
		})
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid groupID",
		})
		return
	}

	exists, err := chats.ChatExists(app.DB, groupID)
	if err != nil {
		log.Println(err, "here1")
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not verify chat",
		})
		return
	}

	if !exists {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid group",
		})
		return
	}

	userIN, err := chats.UserInGroup(app.DB, userID, groupID)
	if err != nil {
		log.Println(err, "here2")
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not verify group",
		})
		return
	}

	if !userIN {
		helpers.WriteJson(w, http.StatusBadRequest, map[string]any{
			"status":  false,
			"message": "invalid request",
		})
		return
	}

	msgs, err := chats.GetMessages(app.DB, userID, groupID, offset)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err, "here3")
		helpers.WriteJson(w, http.StatusInternalServerError, map[string]any{
			"status":  false,
			"message": "could not get messages",
		})
		return
	}

	helpers.WriteJson(w, http.StatusOK, map[string]any{
		"status": true,
		"data":   msgs,
	})
}

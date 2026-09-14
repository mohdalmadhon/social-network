package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	chatsdb "social/database/chats"
)

type chatMessageInput struct {
	Content string `json:"content"`
}

func (app App) PrivateChats(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"status": false, "message": "authentication required"})
		return
	}

	conversations, err := chatsdb.ListPrivateChats(app.DB, userID)
	if err != nil {
		log.Printf("list private chats: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"status": false, "message": "could not load chats"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": true, "chats": conversations})
}

func (app App) PrivateChatUsers(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"status": false, "message": "authentication required"})
		return
	}

	users, err := chatsdb.ListPrivateCandidates(app.DB, userID)
	if err != nil {
		log.Printf("list private chat users: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"status": false, "message": "could not load chat contacts"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": true, "users": users})
}

func (app App) OpenPrivateChat(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"status": false, "message": "authentication required"})
		return
	}

	var input struct {
		UserID int `json:"userId"`
	}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil || input.UserID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"status": false, "message": "valid userId is required"})
		return
	}

	conversation, err := chatsdb.OpenPrivateChat(app.DB, userID, input.UserID)
	if err != nil {
		writeChatError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": true, "chat": conversation})
}

func (app App) PrivateChatMessages(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticatedUserID(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"status": false, "message": "authentication required"})
		return
	}
	chatID, err := parseChatPathID(r, "chatID")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"status": false, "message": "invalid chat id"})
		return
	}

	if r.Method == http.MethodGet {
		messages, listErr := chatsdb.ListPrivateMessages(app.DB, chatID, userID)
		if listErr != nil {
			writeChatError(w, listErr)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": true, "messages": messages})
		return
	}

	input, ok := decodeMessageInput(w, r)
	if !ok {
		return
	}
	message, err := chatsdb.SendPrivateMessage(app.DB, chatID, userID, input.Content)
	if err != nil {
		writeChatError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"status": true, "message": message})
}

func (app App) GroupChatMessages(w http.ResponseWriter, r *http.Request) {
	userID, groupID, ok := groupRequestIdentity(w, r)
	if !ok {
		return
	}

	if r.Method == http.MethodGet {
		messages, err := chatsdb.ListGroupMessages(app.DB, groupID, userID)
		if err != nil {
			writeChatError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": true, "messages": messages})
		return
	}

	input, ok := decodeMessageInput(w, r)
	if !ok {
		return
	}
	message, err := chatsdb.SendGroupMessage(app.DB, groupID, userID, input.Content)
	if err != nil {
		writeChatError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"status": true, "message": message})
}

func decodeMessageInput(w http.ResponseWriter, r *http.Request) (chatMessageInput, bool) {
	var input chatMessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"status": false, "message": "invalid request body"})
		return chatMessageInput{}, false
	}
	return input, true
}

func parseChatPathID(r *http.Request, key string) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue(key), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func writeChatError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, chatsdb.ErrInvalidMessage):
		writeJSON(w, http.StatusBadRequest, map[string]any{"status": false, "message": "message content is required"})
	case errors.Is(err, chatsdb.ErrMessageTooLong):
		writeJSON(w, http.StatusBadRequest, map[string]any{"status": false, "message": "message must be 2000 characters or fewer"})
	case errors.Is(err, chatsdb.ErrSelfChat):
		writeJSON(w, http.StatusBadRequest, map[string]any{"status": false, "message": "you cannot chat with yourself"})
	case errors.Is(err, chatsdb.ErrNoFollowRelation):
		writeJSON(w, http.StatusForbidden, map[string]any{"status": false, "message": "an accepted follow relationship is required"})
	case errors.Is(err, chatsdb.ErrForbidden):
		writeJSON(w, http.StatusForbidden, map[string]any{"status": false, "message": "you do not have access to this chat"})
	case errors.Is(err, chatsdb.ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, map[string]any{"status": false, "message": "user not found"})
	case errors.Is(err, chatsdb.ErrChatNotFound):
		writeJSON(w, http.StatusNotFound, map[string]any{"status": false, "message": "chat not found"})
	case errors.Is(err, chatsdb.ErrGroupNotFound):
		writeJSON(w, http.StatusNotFound, map[string]any{"status": false, "message": "group not found"})
	default:
		log.Printf("chat operation: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"status": false, "message": "could not complete chat operation"})
	}
}

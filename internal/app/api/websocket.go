package api

import (
	"encoding/json"
	"io"
	"log"
	"social/database/chats"
	"social/database/users"
	"social/internal/models"

	"golang.org/x/net/websocket"
)

func (app *App) HandleWS(ws *websocket.Conn) {
	userID, ok := ws.Request().Context().Value("userID").(int)
	if !ok {
		// do somehting
		log.Println("err here")
		return
	}
	log.Println("incoming user")
	app.register(userID, ws)
	defer app.unregister(userID)

	app.readLoop(userID, ws)
}

func (app *App) readLoop(userID int, ws *websocket.Conn) {
	buff := make([]byte, 4096)

	for {
		n, err := ws.Read(buff)

		if err != nil {
			if err == io.EOF {
				break
			}

			log.Println(err)
			break
		}
		log.Println("incoming message")
		var payload models.WSPayload

		if err := json.Unmarshal(buff[:n], &payload); err != nil {
			log.Println(err)
			continue
		}

		switch payload.Type {
		case "privateMessage":
			app.handleMessage(userID, payload.Data)

		case "notification":
			// app.handleNotification(payload.Data)

		default:
			log.Println("unknown websocket type:", payload.Type)
		}
	}
}

func (app *App) handleMessage(userID int, data json.RawMessage) {
	var msg models.IncomingMessage

	if err := json.Unmarshal(data, &msg); err != nil {
		log.Println(err)
		return
	}

	if msg.Content == "" {
		return
	}

	groupID := msg.GroupID
	if groupID <= 0 {
		existingGroupID, err := chats.HasPrivateChat(app.DB, userID, msg.UserID)
		if err != nil {
			return
		}

		if existingGroupID != -1 {
			groupID = existingGroupID
		} else {

			groupID, err = chats.MakePrivateChat(app.DB, userID, msg.UserID)
			if err != nil {
				return
			}
		}
	}

	err := chats.AddMessages(app.DB, msg.Content, userID, groupID)
	if err != nil {
		log.Println(err)
		return
		// do something
	}
	sender, err := users.GetUserSimpleData(app.DB, userID)
	if err != nil {
		log.Println(err)
		return
	}
	message := models.Message{
		Content: msg.Content,
		Sender:  models.UserRegistration{
			ID: userID,
			FirstName: sender.FirstName,
			LastName: sender.LastName,
			Avatar: sender.Avatar,
		},
		GroupID: groupID,
	}
	app.sendToUsers(message, groupID, userID)
}

func (app *App) sendToUsers(msg models.Message, groupID, userID int) {
	ids, err := chats.GetGroupMembersIds(app.DB, groupID)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := json.Marshal(map[string]any{
		"type": "message",
		"data": msg,
	})
	if err != nil {
		log.Println(err)
		return
	}

	for _, id := range ids {
		if id == userID {
			continue
		}

		app.H.Mu.RLock()
		client, ok := app.H.Conn[id]
		app.H.Mu.RUnlock()

		if !ok {
			log.Println("not ok err")
			continue
		}
		
		if _, err := client.Write(response); err != nil {
			log.Println("websocket write error:", err)
		}
	}
}

func (app *App) unregister(userID int) {
	app.H.Mu.Lock()
	defer app.H.Mu.Unlock()
	delete(app.H.Conn, userID)
}

func (app *App) register(userID int, ws *websocket.Conn) {
	app.H.Mu.Lock()
	app.H.Conn[userID] = ws
	app.H.Mu.Unlock()
}

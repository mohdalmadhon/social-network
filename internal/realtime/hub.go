package realtime

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8 * 1024
	sendBufferSize = 64
)

// Hub owns all live WebSocket clients. A user may have several clients open
// at once, such as multiple tabs or devices.
type Hub struct {
	mu      sync.RWMutex
	clients map[int]map[*Client]struct{}
}

type Client struct {
	UserID int
	conn   *websocket.Conn
	hub    *Hub
	send   chan []byte
	done   chan struct{}
	once   sync.Once
}

func NewHub() *Hub {
	return &Hub{clients: make(map[int]map[*Client]struct{})}
}

func (h *Hub) Register(userID int, conn *websocket.Conn) *Client {
	client := &Client{
		UserID: userID,
		conn:   conn,
		hub:    h,
		send:   make(chan []byte, sendBufferSize),
		done:   make(chan struct{}),
	}

	h.mu.Lock()
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]struct{})
	}
	h.clients[userID][client] = struct{}{}
	h.mu.Unlock()

	go client.writePump()
	return client
}

func (h *Hub) Unregister(client *Client) {
	if client == nil {
		return
	}
	client.once.Do(func() {
		h.mu.Lock()
		if userClients := h.clients[client.UserID]; userClients != nil {
			if _, exists := userClients[client]; exists {
				delete(userClients, client)
				close(client.done)
			}
			if len(userClients) == 0 {
				delete(h.clients, client.UserID)
			}
		}
		h.mu.Unlock()
	})
}

func (h *Hub) SendToUser(userID int, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	h.send(userID, data)
	return nil
}

func (h *Hub) SendToUsers(userIDs []int, payloadForUser func(int) any) error {
	for _, userID := range userIDs {
		if err := h.SendToUser(userID, payloadForUser(userID)); err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) ConnectionCount(userID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID])
}

func (h *Hub) DisconnectUser(userID int) {
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients[userID]))
	for client := range h.clients[userID] {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		h.Unregister(client)
		_ = client.conn.Close()
	}
}

func (h *Hub) send(userID int, data []byte) {
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients[userID]))
	for client := range h.clients[userID] {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		select {
		case <-client.done:
			continue
		default:
		}
		select {
		case client.send <- data:
		default:
			h.Unregister(client)
			_ = client.conn.Close()
		}
	}
}

func (c *Client) ReadPump(handle func([]byte)) {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		messageType, payload, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			continue
		}
		handle(payload)
	}
}

func (c *Client) Send(payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	select {
	case <-c.done:
		return nil
	default:
	}
	select {
	case c.send <- data:
	default:
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}
	return nil
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				c.hub.Unregister(c)
				return
			}
		case <-c.done:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.hub.Unregister(c)
				return
			}
		}
	}
}

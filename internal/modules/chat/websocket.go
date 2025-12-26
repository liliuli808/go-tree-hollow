package chat

import (
	"encoding/json"
	"go-tree-hollow/internal/models"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type      string          `json:"type"` // message, typing, read, online
	To        uint            `json:"to,omitempty"`
	From      uint            `json:"from,omitempty"`
	Content   string          `json:"content,omitempty"`
	MessageID uint            `json:"message_id,omitempty"`
	Message   *models.Message `json:"message,omitempty"`
}

// Client represents a WebSocket client connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userID uint
	send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[uint]*Client // userID -> client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *WebSocketMessage
	service    Service
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub(service Service) *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *WebSocketMessage),
		service:    service,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("User %d connected to WebSocket", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("User %d disconnected from WebSocket", client.userID)

		case msg := <-h.broadcast:
			h.handleBroadcast(msg)
		}
	}
}

// handleBroadcast handles broadcasting messages to specific users
func (h *Hub) handleBroadcast(msg *WebSocketMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.mu.RLock()
	if client, ok := h.clients[msg.To]; ok {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, msg.To)
		}
	}
	h.mu.RUnlock()
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID uint, msg *WebSocketMessage) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return // User not connected
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case client.send <- data:
	default:
		log.Printf("Send buffer full for user %d", userID)
	}
}

// IsUserOnline checks if a user is connected
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// HandleConnection upgrades HTTP connection to WebSocket
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request, userID uint) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
	}

	h.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024) // 512KB max message size
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMsg WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		c.handleMessage(&wsMsg)
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(msg *WebSocketMessage) {
	msg.From = c.userID

	switch msg.Type {
	case "message":
		// Save message to database
		savedMsg, err := c.hub.service.SendMessage(c.userID, msg.To, msg.Content)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			return
		}

		// Send back to sender (confirmation)
		msg.Message = savedMsg
		c.hub.SendToUser(c.userID, msg)

		// Send to receiver
		c.hub.SendToUser(msg.To, msg)

	case "typing":
		// Forward typing indicator to receiver
		c.hub.SendToUser(msg.To, msg)

	case "read":
		// Mark message as read
		if err := c.hub.service.MarkAsRead(msg.MessageID, c.userID); err != nil {
			log.Printf("Error marking message as read: %v", err)
			return
		}
		// Notify sender that message was read
		c.hub.SendToUser(msg.To, msg)
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

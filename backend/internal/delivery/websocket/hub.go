package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Message types
const (
	MsgTypeHeartbeat         = "heartbeat"
	MsgTypeViolation         = "violation"
	MsgTypeParticipantJoin   = "participant_join"
	MsgTypeParticipantLeave  = "participant_leave"
	MsgTypeParticipantFinish = "participant_finish"
	MsgTypeSessionEnd        = "session_end"
	MsgTypeCameraFrame       = "camera_frame"
)

type Message struct {
	Type          string          `json:"type"`
	SessionID     string          `json:"session_id,omitempty"`
	ParticipantID string          `json:"participant_id,omitempty"`
	UserName      string          `json:"user_name,omitempty"`
	Data          json.RawMessage `json:"data,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
}

type Client struct {
	ID            string
	SessionID     string
	Role          string // "hr" or "participant"
	ParticipantID string
	UserName      string
	conn          *websocket.Conn
	send          chan []byte
	hub           *Hub
}

type Hub struct {
	mu       sync.RWMutex
	sessions map[string]map[*Client]bool // sessionID -> set of clients
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.sessions[client.SessionID] == nil {
		h.sessions[client.SessionID] = make(map[*Client]bool)
	}
	h.sessions[client.SessionID][client] = true
	log.Printf("WS Client registered: %s in session %s", client.ID, client.SessionID)
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.sessions[client.SessionID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.sessions, client.SessionID)
		}
	}
	close(client.send)
	log.Printf("WS Client unregistered: %s", client.ID)
}

// BroadcastToSession sends message to all HR clients in the session
func (h *Hub) BroadcastToSession(sessionID string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, _ := json.Marshal(msg)
	for client := range h.sessions[sessionID] {
		if client.Role == "hr" {
			select {
			case client.send <- data:
			default:
				// drop if buffer full
			}
		}
	}
}

// BroadcastToSessionAll sends message to all clients (HR and participants)
func (h *Hub) BroadcastToSessionAll(sessionID string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, _ := json.Marshal(msg)
	for client := range h.sessions[sessionID] {
		select {
		case client.send <- data:
		default:
			// drop if buffer full
		}
	}
}

// NotifyParticipant sends a message to a specific participant
func (h *Hub) NotifyParticipant(sessionID, participantID string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, _ := json.Marshal(msg)
	for client := range h.sessions[sessionID] {
		if client.ParticipantID == participantID {
			select {
			case client.send <- data:
			default:
			}
		}
	}
}

func (h *Hub) StartClient(client *Client) {
	go client.writePump()
	go client.readPump(h)
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
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

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024) // 512 KB — large enough for base64 JPEG camera frames
	c.conn.SetReadDeadline(time.Now().Add(120 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		return nil
	})

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			continue
		}

		msg.SessionID = c.SessionID
		msg.ParticipantID = c.ParticipantID
		msg.UserName = c.UserName
		msg.Timestamp = time.Now()

		// Forward participant events to HR dashboard
		if c.Role == "participant" {
			h.BroadcastToSession(c.SessionID, msg)
		}
	}
}

// NewClient creates a new WebSocket client
func NewClient(id, sessionID, role, participantID, userName string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:            id,
		SessionID:     sessionID,
		Role:          role,
		ParticipantID: participantID,
		UserName:      userName,
		conn:          conn,
		send:          make(chan []byte, 64),
		hub:           hub,
	}
}

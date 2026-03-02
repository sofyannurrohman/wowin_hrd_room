package handlers

import (
	"net/http"

	ws "hrd_room/backend/internal/delivery/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (configure for production)
	},
}

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(hub *ws.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// GET /api/ws/:sessionId
// Query params: role=hr|participant, participant_id=UUID, user_name=string
func (h *WSHandler) Handle(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	role := c.Query("role")
	if role == "" {
		role = "participant"
	}

	participantID := c.Query("participant_id")
	userName := c.Query("user_name")
	if userName == "" {
		userNameI, _ := c.Get("user_name")
		if userNameI != nil {
			userName = userNameI.(string)
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(
		uuid.New().String(),
		sessionID,
		role,
		participantID,
		userName,
		conn,
		h.hub,
	)

	h.hub.Register(client)
	h.hub.StartClient(client)
}

package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	ws "hrd_room/backend/internal/delivery/websocket"
	"hrd_room/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionHandler struct {
	sessionUC *usecase.SessionUseCase
	wsHub     *ws.Hub
	uploadDir string
}

func NewSessionHandler(sessionUC *usecase.SessionUseCase, hub *ws.Hub, uploadDir string) *SessionHandler {
	return &SessionHandler{sessionUC: sessionUC, wsHub: hub, uploadDir: uploadDir}
}

func (h *SessionHandler) Create(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req usecase.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	session, err := h.sessionUC.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusCreated, session)
}

func (h *SessionHandler) List(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	role, _ := c.Get("user_role")

	sessions, err := h.sessionUC.List(c.Request.Context(), role.(string), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	session, err := h.sessionUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	var req usecase.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	session, err := h.sessionUC.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	// Notify all active participants via WebSocket BEFORE deleting,
	// so their ExamPage auto-submits answers and redirects to /exam-finished.
	if h.wsHub != nil {
		h.wsHub.BroadcastToSessionAll(id.String(), ws.Message{
			Type:      ws.MsgTypeSessionEnd,
			SessionID: id.String(),
			Timestamp: time.Now(),
		})
	}

	if err := h.sessionUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "session deleted"})
}

// PUT /api/sessions/:id/lock
func (h *SessionHandler) Lock(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	var req struct {
		Locked bool `json:"locked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	if err := h.sessionUC.Lock(c.Request.Context(), id, req.Locked); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}

	// Broadcast session_end message to all clients when locking the session
	if req.Locked && h.wsHub != nil {
		h.wsHub.BroadcastToSessionAll(id.String(), ws.Message{
			Type:      ws.MsgTypeSessionEnd,
			SessionID: id.String(),
			Timestamp: time.Now(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "session updated"})
}

func (h *SessionHandler) GenerateTokens(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	var req usecase.GenerateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	tokens, err := h.sessionUC.GenerateTokens(c.Request.Context(), sessionID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})
}

func (h *SessionHandler) ListTokens(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	tokens, err := h.sessionUC.ListTokens(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// Question handlers (within session context)

type QuestionHandler struct {
	questionRepo interface {
		Create(ctx interface{}, q interface{}) error
	}
	uploadDir string
}

func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' {
			return r
		}
		return -1
	}, name)
	return name
}

func generateUploadPath(uploadDir, prefix, ext string) string {
	filename := sanitizeFilename(prefix) + "_" + time.Now().Format("20060102150405") + ext
	return filepath.Join(uploadDir, "questions", filename)
}

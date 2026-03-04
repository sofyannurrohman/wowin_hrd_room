package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	ws "hrd_room/backend/internal/delivery/websocket"
	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"
	"hrd_room/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExamHandler struct {
	examUC    *usecase.ExamUseCase
	sessionUC *usecase.SessionUseCase
	qRepo     *repository.QuestionRepository
	vRepo     *repository.ViolationRepository
	wsHub     *ws.Hub
	uploadDir string
}

func NewExamHandler(examUC *usecase.ExamUseCase, sessionUC *usecase.SessionUseCase, qRepo *repository.QuestionRepository, vRepo *repository.ViolationRepository, hub *ws.Hub, uploadDir string) *ExamHandler {
	return &ExamHandler{examUC: examUC, sessionUC: sessionUC, qRepo: qRepo, vRepo: vRepo, wsHub: hub, uploadDir: uploadDir}
}

// POST /api/exam/join
func (h *ExamHandler) Join(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Age      int    `json:"age" binding:"required"`
		Position string `json:"position" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, session, err := h.sessionUC.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participant, err := h.examUC.Join(c.Request.Context(), req.Name, req.Email, req.Age, req.Position, token, session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compute remaining time until session end for timer initialization
	endTime := session.Schedule.Add(time.Duration(session.DurationMinutes) * time.Minute)
	remaining := int(endTime.Sub(time.Now()).Seconds())
	if remaining < 0 {
		remaining = 0
	}

	// Notify HR dashboards that a participant has joined
	if h.wsHub != nil {
		h.wsHub.BroadcastToSession(session.ID.String(), ws.Message{
			Type:          ws.MsgTypeParticipantJoin,
			SessionID:     session.ID.String(),
			ParticipantID: participant.ID.String(),
			UserName:      req.Name,
			Timestamp:     time.Now(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"participant_id": participant.ID,
		"session_id":     session.ID,
		"session_name":   session.Name,
		"duration":       session.DurationMinutes,
		"time_remaining": remaining,
	})
}

// GET /api/exam/:sessionId/modules & /api/sessions/:id/modules
func (h *ExamHandler) GetModules(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	if sessionIDStr == "" {
		sessionIDStr = c.Param("id")
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	modules, err := h.examUC.GetModulesForParticipant(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if modules == nil {
		modules = []domain.Module{}
	}
	c.JSON(http.StatusOK, modules)
}

// GET /api/exam/:sessionId/modules/:moduleId/questions
func (h *ExamHandler) GetQuestionsForModule(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("sessionId"))
	moduleID, err2 := uuid.Parse(c.Param("moduleId"))
	if err != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ids"})
		return
	}

	questions, err := h.examUC.GetQuestionsForModule(c.Request.Context(), sessionID, moduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// POST /api/exam/:sessionId/answers
func (h *ExamHandler) SubmitAnswers(c *gin.Context) {
	var req usecase.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant_id"})
		return
	}

	if err := h.examUC.SaveAnswers(c.Request.Context(), participantID, req.Answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Auto-grade immediately on submit
	result, err := h.examUC.AutoGrade(c.Request.Context(), participantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "jawaban berhasil disubmit",
		"result_id":      result.ID,
		"grading_status": result.GradingStatus,
	})
}

// POST /api/exam/:sessionId/answers/autosave
func (h *ExamHandler) AutoSaveAnswers(c *gin.Context) {
	var req usecase.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant_id"})
		return
	}

	if err := h.examUC.SaveAnswers(c.Request.Context(), participantID, req.Answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "jawaban berhasil disimpan otomatis",
	})
}

// GET /api/exam/answers/:participantId
func (h *ExamHandler) GetParticipantAnswersPublic(c *gin.Context) {
	participantID, err := uuid.Parse(c.Param("participantId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant id"})
		return
	}
	answers, err := h.examUC.GetParticipantAnswers(c.Request.Context(), participantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, answers)
}

// GET /api/sessions/:id/results
func (h *ExamHandler) GetResults(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	results, err := h.examUC.GetResults(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// GET /api/results/:participantId/answers
func (h *ExamHandler) GetParticipantAnswers(c *gin.Context) {
	participantID, err := uuid.Parse(c.Param("participantId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant id"})
		return
	}
	answers, err := h.examUC.GetParticipantAnswers(c.Request.Context(), participantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, answers)
}

// PUT /api/results/:id/review
func (h *ExamHandler) HRReview(c *gin.Context) {
	var req usecase.HRReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.examUC.HRReview(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "review berhasil disimpan"})
}

// POST /api/results/:id/finalize
func (h *ExamHandler) FinalizeScore(c *gin.Context) {
	resultID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid result id"})
		return
	}
	if err := h.examUC.FinalizeScore(c.Request.Context(), resultID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "nilai berhasil difinalisasi"})
}

// GET /api/sessions/:id/participants
func (h *ExamHandler) GetParticipants(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	participants, err := h.examUC.GetParticipants(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, participants)
}

// POST /api/violations
func (h *ExamHandler) ReportViolation(c *gin.Context) {
	var req struct {
		ParticipantID string `json:"participant_id" binding:"required"`
		SessionID     string `json:"session_id" binding:"required"`
		ViolationType string `json:"violation_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant_id"})
		return
	}

	v := &domain.Violation{
		ID:            uuid.New(),
		ParticipantID: participantID,
		ViolationType: req.ViolationType,
	}

	if err := h.vRepo.Create(c.Request.Context(), v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Broadcast violation event to HR dashboards over WebSocket
	if h.wsHub != nil {
		payload, _ := json.Marshal(map[string]string{
			"violation_type": req.ViolationType,
			"participant_id": req.ParticipantID,
		})
		h.wsHub.BroadcastToSession(req.SessionID, ws.Message{
			Type:          ws.MsgTypeViolation,
			SessionID:     req.SessionID,
			ParticipantID: req.ParticipantID,
			Data:          payload,
			Timestamp:     time.Now(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{"id": v.ID, "message": "pelanggaran dicatat"})
}

// GET /api/sessions/:id/violations
func (h *ExamHandler) GetViolations(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	violations, err := h.vRepo.ListBySession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, violations)
}

// ─── Question CRUD Handlers ──────────────────────────────────────────────────

type QuestionHTTPHandler struct {
	qRepo     *repository.QuestionRepository
	uploadDir string
}

func NewQuestionHTTPHandler(qRepo *repository.QuestionRepository, uploadDir string) *QuestionHTTPHandler {
	return &QuestionHTTPHandler{qRepo: qRepo, uploadDir: uploadDir}
}

// POST /api/questions
func (h *QuestionHTTPHandler) Create(c *gin.Context) {
	// Handle multipart for file uploads
	moduleIDStr := c.PostForm("module_id")
	qType := c.PostForm("type")
	content := c.PostForm("content")

	if moduleIDStr == "" || qType == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "module_id, type, content required"})
		return
	}

	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module_id"})
		return
	}

	q := &domain.Question{
		ID:                   uuid.New(),
		ModuleID:             moduleID,
		Type:                 qType,
		Content:              content,
		Weight:               1.0,
		RequiresManualReview: qType == domain.QuestionTypePsychological || qType == domain.QuestionTypeShortAnswer,
	}

	// Handle image upload
	file, header, err := c.Request.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !allowed[strings.ToLower(ext)] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format gambar tidak didukung"})
			return
		}

		imgDir := filepath.Join(h.uploadDir, "questions")
		if err := os.MkdirAll(imgDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat direktori upload"})
			return
		}

		filename := fmt.Sprintf("%s_%s%s", q.ID.String(), time.Now().Format("20060102150405"), ext)
		dstPath := filepath.Join(imgDir, filename)

		if err := c.SaveUploadedFile(header, dstPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
			return
		}

		imgURL := "/uploads/questions/" + filename
		q.ImageURL = &imgURL
	}

	if err := h.qRepo.Create(c.Request.Context(), q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle options from form (JSON string)
	optionsJSON := c.PostForm("options")
	if optionsJSON != "" {
		var options []struct {
			Content   string `json:"content"`
			IsCorrect bool   `json:"is_correct"`
		}
		if err := (&gin.Context{}).ShouldBindJSON(&options); err == nil {
			for _, opt := range options {
				o := &domain.QuestionOption{
					ID:         uuid.New(),
					QuestionID: q.ID,
					Content:    opt.Content,
					IsCorrect:  opt.IsCorrect,
				}
				_ = h.qRepo.CreateOption(c.Request.Context(), o)
			}
		}
	}

	q, _ = h.qRepo.FindByID(c.Request.Context(), q.ID)
	c.JSON(http.StatusCreated, q)
}

// GET /api/modules/:id/questions (for HR)
func (h *QuestionHTTPHandler) ListByModule(c *gin.Context) {
	moduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module id"})
		return
	}
	questions, err := h.qRepo.ListByModule(c.Request.Context(), moduleID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// GET /api/questions
func (h *QuestionHTTPHandler) ListAll(c *gin.Context) {
	questions, err := h.qRepo.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// GET /api/questions/:id
func (h *QuestionHTTPHandler) GetByID(c *gin.Context) {
	qID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}
	q, err := h.qRepo.FindByID(c.Request.Context(), qID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}
	c.JSON(http.StatusOK, q)
}

// PUT /api/questions/:id
func (h *QuestionHTTPHandler) Update(c *gin.Context) {
	qID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}

	q, err := h.qRepo.FindByID(c.Request.Context(), qID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}

	if t := c.PostForm("type"); t != "" {
		q.Type = t
	}
	if content := c.PostForm("content"); content != "" {
		q.Content = content
	}

	// Handle new image upload
	file, header, err := c.Request.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		imgDir := filepath.Join(h.uploadDir, "questions")
		os.MkdirAll(imgDir, 0755)
		filename := fmt.Sprintf("%s_%s%s", q.ID.String(), time.Now().Format("20060102150405"), ext)
		dstPath := filepath.Join(imgDir, filename)
		if err := c.SaveUploadedFile(header, dstPath); err == nil {
			imgURL := "/uploads/questions/" + filename
			q.ImageURL = &imgURL
		}
	}

	if err := h.qRepo.Update(c.Request.Context(), q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update options if provided
	optionsJSON := c.PostForm("options")
	if optionsJSON != "" {
		var options []struct {
			Content   string `json:"content"`
			IsCorrect bool   `json:"is_correct"`
		}
		// Manually unmarshal since Context ShouldBindJSON requires application/json and this is a form field
		if err := json.Unmarshal([]byte(optionsJSON), &options); err == nil {
			_ = h.qRepo.DeleteOptions(c.Request.Context(), q.ID)
			for _, opt := range options {
				o := &domain.QuestionOption{
					ID:         uuid.New(),
					QuestionID: q.ID,
					Content:    opt.Content,
					IsCorrect:  opt.IsCorrect,
				}
				_ = h.qRepo.CreateOption(c.Request.Context(), o)
			}
		}
	}

	q, _ = h.qRepo.FindByID(c.Request.Context(), qID)
	c.JSON(http.StatusOK, q)
}

// DELETE /api/questions/:id
func (h *QuestionHTTPHandler) Delete(c *gin.Context) {
	qID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}
	if err := h.qRepo.Delete(c.Request.Context(), qID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "soal berhasil dihapus"})
}

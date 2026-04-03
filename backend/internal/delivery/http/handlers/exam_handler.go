package handlers

import (
	"encoding/csv"
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
		Name     string `json:"name"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
		Position string `json:"position"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	token, session, err := h.sessionUC.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	// Frictionless Join: If token is bound to an email, bypass form data
	if token.BoundEmail != nil && *token.BoundEmail != "" {
		req.Email = *token.BoundEmail
		// Name, Age, Position can be empty here; ExamUseCase will look up by Email.
	} else {
		// If no bound email, form data is strictly required
		if req.Name == "" || req.Email == "" || req.Position == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nama, email, dan posisi wajib diisi jika tidak menggunakan token khusus"})
			return
		}
	}

	participant, err := h.examUC.Join(c.Request.Context(), req.Name, req.Email, req.Age, req.Position, token, session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}

	// Compute remaining time until session end for timer initialization
	endTime := session.Schedule.Add(time.Duration(session.DurationMinutes) * time.Minute)
	remaining := int(time.Until(endTime).Seconds())
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// POST /api/exam/:sessionId/answers
func (h *ExamHandler) SubmitAnswers(c *gin.Context) {
	var req usecase.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant_id"})
		return
	}

	if err := h.examUC.SaveAnswers(c.Request.Context(), participantID, req.Answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}

	// Auto-grade immediately on submit
	result, err := h.examUC.AutoGrade(c.Request.Context(), participantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	participantID, err := uuid.Parse(req.ParticipantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant_id"})
		return
	}

	if err := h.examUC.SaveAnswers(c.Request.Context(), participantID, req.Answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, answers)
}

// GET /api/exam/participant/:participantId/status
// Public endpoint used by the frontend router guard to check if a stored participantId
// is still active before resuming an exam session (prevents blank page on stale localStorage).
func (h *ExamHandler) GetParticipantStatus(c *gin.Context) {
	participantID, err := uuid.Parse(c.Param("participantId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant id"})
		return
	}
	participant, err := h.examUC.GetParticipantByID(c.Request.Context(), participantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "participant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": participant.Status})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, answers)
}

// PUT /api/results/:id/review
func (h *ExamHandler) HRReview(c *gin.Context) {
	var req usecase.HRReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}
	if err := h.examUC.HRReview(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, violations)
}

// ─── Question CRUD Handlers ──────────────────────────────────────────────────

type QuestionHTTPHandler struct {
	qRepo     *repository.QuestionRepository
	mRepo     *repository.ModuleRepository
	uploadDir string
}

func NewQuestionHTTPHandler(qRepo *repository.QuestionRepository, mRepo *repository.ModuleRepository, uploadDir string) *QuestionHTTPHandler {
	return &QuestionHTTPHandler{qRepo: qRepo, mRepo: mRepo, uploadDir: uploadDir}
}

// POST /api/questions
func (h *QuestionHTTPHandler) Create(c *gin.Context) {
	// Handle multipart for file uploads
	moduleIDStr := c.PostForm("module_id")
	qType := normalizeQuestionType(c.PostForm("type"))
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
	weightStr := c.PostForm("weight")
	weight := 1.0
	if weightStr != "" {
		fmt.Sscanf(weightStr, "%f", &weight)
	}

	module, err := h.mRepo.FindByID(c.Request.Context(), moduleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "modul tidak ditemukan"})
		return
	}

	currentWeight, err := h.qRepo.GetTotalWeightByModule(c.Request.Context(), moduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung bobot"})
		return
	}

	if currentWeight+weight > module.TotalWeight {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Bobot pertanyaan melebihi batas maksimal modul (Sisa: %.2f)", module.TotalWeight-currentWeight)})
		return
	}

	q := &domain.Question{
		ID:                   uuid.New(),
		ModuleID:             moduleID,
		Type:                 qType,
		Content:              content,
		Weight:               weight,
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

		if header.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar maksimal 2MB"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}

	// Handle options from form (JSON string)
	optionsJSON := c.PostForm("options")
	if optionsJSON != "" {
		var options []struct {
			Content   string `json:"content"`
			IsCorrect bool   `json:"is_correct"`
		}
		if err := json.Unmarshal([]byte(optionsJSON), &options); err == nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// GET /api/questions
func (h *QuestionHTTPHandler) ListAll(c *gin.Context) {
	questions, err := h.qRepo.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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
		q.Type = normalizeQuestionType(t)
	}
	if content := c.PostForm("content"); content != "" {
		q.Content = content
	}
	if weightStr := c.PostForm("weight"); weightStr != "" {
		var newWeight float64
		if _, err := fmt.Sscanf(weightStr, "%f", &newWeight); err == nil {
			if newWeight != q.Weight {
				module, err := h.mRepo.FindByID(c.Request.Context(), q.ModuleID)
				if err == nil {
					currentTotal, _ := h.qRepo.GetTotalWeightByModule(c.Request.Context(), q.ModuleID)
					weightDifference := newWeight - q.Weight
					if currentTotal+weightDifference > module.TotalWeight {
						c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Bobot pertanyaan melebihi batas maksimal modul (Sisa: %.2f)", module.TotalWeight-currentTotal)})
						return
					}
				}
				q.Weight = newWeight
			}
		}
	}

	// Handle new image upload
	file, header, err := c.Request.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)

		if header.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar maksimal 2MB"})
			return
		}

		imgDir := filepath.Join(h.uploadDir, "questions")
		os.MkdirAll(imgDir, 0755)
		filename := fmt.Sprintf("%s_%s%s", q.ID.String(), time.Now().Format("20060102150405"), ext)
		dstPath := filepath.Join(imgDir, filename)
		if err := c.SaveUploadedFile(header, dstPath); err == nil {
			// Delete old image if it exists
			if q.ImageURL != nil && *q.ImageURL != "" {
				oldRelPath := strings.TrimPrefix(*q.ImageURL, "/uploads")
				oldFullPath := filepath.Join(h.uploadDir, oldRelPath)
				os.Remove(oldFullPath)
			}
			imgURL := "/uploads/questions/" + filename
			q.ImageURL = &imgURL
		}
	}

	if err := h.qRepo.Update(c.Request.Context(), q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
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

// POST /api/questions/import
func (h *QuestionHTTPHandler) ImportCSV(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file CSV tidak ditemukan"})
		return
	}
	defer file.Close()

	if fileHeader.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file CSV maksimal 2MB"})
		return
	}

	reader := csv.NewReader(file)
	// Optionally handle different delimiters, though default is comma
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format CSV tidak valid"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV kosong atau tidak memiliki baris header/data"})
		return
	}

	successCount := 0
	errorCount := 0

	moduleCache := make(map[uuid.UUID]*domain.Module)
	moduleCurrentWeights := make(map[uuid.UUID]float64)

	// Skip header (row 0)
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < 3 { // Min required columns: content, type, module_id
			errorCount++
			continue
		}

		content := strings.TrimSpace(row[0])
		if content == "" {
			errorCount++
			continue
		}

		qType := normalizeQuestionType(row[1])

		moduleIDStr := strings.TrimSpace(row[2])
		var moduleID uuid.UUID
		if moduleIDStr != "" {
			parsedID, err := uuid.Parse(moduleIDStr)
			if err == nil {
				moduleID = parsedID
			}
		}

		weight := 1.0
		if len(row) > 9 {
			wStr := strings.TrimSpace(row[9])
			if wStr != "" {
				fmt.Sscanf(wStr, "%f", &weight)
			}
		}

		// Validate module and its weight limit
		if moduleID == uuid.Nil {
			errorCount++
			continue
		}

		mod, exists := moduleCache[moduleID]
		if !exists {
			var err error
			mod, err = h.mRepo.FindByID(c.Request.Context(), moduleID)
			if err != nil {
				errorCount++
				continue
			}
			moduleCache[moduleID] = mod

			currW, err := h.qRepo.GetTotalWeightByModule(c.Request.Context(), moduleID)
			if err != nil {
				errorCount++
				continue
			}
			moduleCurrentWeights[moduleID] = currW
		}

		if moduleCurrentWeights[moduleID]+weight > mod.TotalWeight {
			// Skip this question as it exceeds module weight limit
			errorCount++
			continue
		}

		q := &domain.Question{
			ID:                   uuid.New(),
			Type:                 qType,
			Content:              content,
			ModuleID:             moduleID,
			Weight:               weight,
			RequiresManualReview: qType == domain.QuestionTypePsychological || qType == domain.QuestionTypeShortAnswer,
		}

		if err := h.qRepo.Create(c.Request.Context(), q); err != nil {
			errorCount++
			continue
		}
		moduleCurrentWeights[moduleID] += weight

		// Handle Options for multiple choice/select/true-false/short-answer
		if qType == domain.QuestionTypeMultipleChoice || qType == domain.QuestionTypeTrueFalse || qType == domain.QuestionTypeShortAnswer {
			var correctAnswers []string
			if len(row) > 8 {
				correctAnsText := strings.TrimSpace(row[8])
				if qType == domain.QuestionTypeMultipleChoice {
					// For MCQ, we expect labels like A, B, C or A,B
					correctAnswers = strings.Split(strings.ToUpper(correctAnsText), ",")
					for j, ca := range correctAnswers {
						correctAnswers[j] = strings.TrimSpace(ca)
					}
				} else {
					// For other types, the text itself is the correct answer
					correctAnswers = []string{correctAnsText}
				}
			}

			// Map A->3, B->4, C->5, D->6, E->7
			optLabels := []string{"A", "B", "C", "D", "E"}
			optionCreated := false
			for j := 0; j < 5; j++ {
				colIdx := 3 + j
				if colIdx < len(row) {
					optContent := strings.TrimSpace(row[colIdx])
					if optContent != "" {
						isCorrect := false
						label := optLabels[j]
						for _, ca := range correctAnswers {
							if ca == label || strings.EqualFold(ca, optContent) {
								isCorrect = true
								break
							}
						}
						o := &domain.QuestionOption{
							ID:         uuid.New(),
							QuestionID: q.ID,
							Content:    optContent,
							IsCorrect:  isCorrect,
						}
						_ = h.qRepo.CreateOption(c.Request.Context(), o)
						optionCreated = true
					}
				}
			}

			// Special case: if no options were provided in columns A-E, 
			// use the correct_answer_text as the single correct option (common for short_answer/true_false)
			if !optionCreated && len(correctAnswers) > 0 && correctAnswers[0] != "" {
				o := &domain.QuestionOption{
					ID:         uuid.New(),
					QuestionID: q.ID,
					Content:    correctAnswers[0],
					IsCorrect:  true,
				}
				_ = h.qRepo.CreateOption(c.Request.Context(), o)
			}
		}

		successCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Import selesai. Berhasil: %d, Gagal: %d", successCount, errorCount),
		"success_count": successCount,
		"error_count":   errorCount,
	})
}

// DELETE /api/questions/:id
func (h *QuestionHTTPHandler) Delete(c *gin.Context) {
	qID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
		return
	}

	// Fetch question to check for image
	q, err := h.qRepo.FindByID(c.Request.Context(), qID)
	if err == nil && q.ImageURL != nil && *q.ImageURL != "" {
		// Attempt to delete image file
		// Image url format is usually /uploads/questions/filename
		relPath := strings.TrimPrefix(*q.ImageURL, "/uploads")
		fullPath := filepath.Join(h.uploadDir, relPath)
		os.Remove(fullPath)
	}

	if err := h.qRepo.Delete(c.Request.Context(), qID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "soal berhasil dihapus"})
}

// POST /api/exam/participant/:participantId/selfie
func (h *ExamHandler) UploadSelfie(c *gin.Context) {
	participantIDStr := c.Param("participantId")
	participantID, err := uuid.Parse(participantIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant_id"})
		return
	}

	file, header, err := c.Request.FormFile("selfie")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file selfie tidak ditemukan"})
		return
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran gambar maksimal 5MB"})
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg" // default if blob doesn't have ext
	}
	
	imgDir := filepath.Join(h.uploadDir, "selfies")
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat direktori upload"})
		return
	}

	filename := fmt.Sprintf("%s_%s%s", participantID.String(), time.Now().Format("20060102150405"), ext)
	dstPath := filepath.Join(imgDir, filename)

	if err := c.SaveUploadedFile(header, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar selfie"})
		return
	}

	var existingSelfiePath string
	p, err := h.examUC.GetParticipantByID(c.Request.Context(), participantID)
	if err == nil && p.KtpSelfieURL != nil && *p.KtpSelfieURL != "" {
		existingSelfiePath = filepath.Join(h.uploadDir, strings.TrimPrefix(*p.KtpSelfieURL, "/uploads"))
	}

	imgURL := "/uploads/selfies/" + filename
	if err := h.examUC.UpdateSelfieUrl(c.Request.Context(), participantID, imgURL); err != nil {
		os.Remove(dstPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}

	if existingSelfiePath != "" {
		os.Remove(existingSelfiePath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "selfie berhasil diunggah", "url": imgURL})
}

func normalizeQuestionType(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "multiple_choice", "multiple choice", "pilihan ganda", "":
		return domain.QuestionTypeMultipleChoice
	case "true_false", "true/false", "benar/salah", "benar salah":
		return domain.QuestionTypeTrueFalse
	case "short_answer", "short answer", "isian singkat":
		return domain.QuestionTypeShortAnswer
	case "psychological", "essay", "esay", "psikologi":
		return domain.QuestionTypePsychological
	default:
		return t
	}
}



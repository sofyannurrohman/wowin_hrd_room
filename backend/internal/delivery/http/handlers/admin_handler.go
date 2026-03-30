package handlers

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userRepo  *repository.UserRepository
	logRepo   *repository.LogRepository
	uploadDir string
}

func NewAdminHandler(userRepo *repository.UserRepository, logRepo *repository.LogRepository, uploadDir string) *AdminHandler {
	return &AdminHandler{userRepo: userRepo, logRepo: logRepo, uploadDir: uploadDir}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.userRepo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), id)
	if err == nil && user.CvURL != nil && *user.CvURL != "" {
		// Attempt to delete CV file
		// CV url format is usually /uploads/cv/filename
		relPath := strings.TrimPrefix(*user.CvURL, "/uploads")
		fullPath := filepath.Join(h.uploadDir, relPath)
		os.Remove(fullPath)
	}

	if err := h.userRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

type UpdateUserRequest struct {
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	RoleID          int     `json:"role_id"`
	Age             *int    `json:"age"`
	AppliedPosition *string `json:"applied_position"`
	Address         *string `json:"address"`
	LastEducation   *string `json:"last_education"`
	WhatsappNumber  *string `json:"whatsapp_number"`
	ExpectedSalary  *string `json:"expected_salary"`
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}
	if req.Age != nil {
		user.Age = req.Age
	}
	if req.AppliedPosition != nil {
		user.AppliedPosition = req.AppliedPosition
	}
	if req.Address != nil {
		user.Address = req.Address
	}
	if req.LastEducation != nil {
		user.LastEducation = req.LastEducation
	}
	if req.WhatsappNumber != nil {
		user.WhatsappNumber = req.WhatsappNumber
	}
	if req.ExpectedSalary != nil {
		user.ExpectedSalary = req.ExpectedSalary
	}

	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AdminHandler) GetLogs(c *gin.Context) {
	logs, err := h.logRepo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// POST /api/participants/import
func (h *AdminHandler) ImportParticipants(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file CSV wajib diunggah"})
		return
	}
	defer file.Close()

	if fileHeader.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file CSV maksimal 2MB"})
		return
	}

	reader := csv.NewReader(file)
	// Read header
	header, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal membaca header CSV"})
		return
	}

	// Map header columns to indices
	colMap := make(map[string]int)
	for i, col := range header {
		colMap[strings.ToLower(strings.TrimSpace(col))] = i
	}

	nameIdx, nameOk := colMap["name"]
	emailIdx, emailOk := colMap["email"]
	if !nameOk || !emailOk {
		// Try indonesian
		nameIdx, nameOk = colMap["nama"]
		emailIdx, emailOk = colMap["email"]
		if !nameOk || !emailOk {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV harus memiliki kolom 'name' atau 'nama', dan 'email'"})
			return
		}
	}

	ageIdx, ageOk := colMap["age"]
	if !ageOk {
		ageIdx, ageOk = colMap["umur"]
	}

	posIdx, posOk := colMap["position"]
	if !posOk {
		posIdx, posOk = colMap["posisi"]
	}

	imported := 0
	skipped := 0
	var skippedEmails []string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) == 0 {
			continue // skip broken rows
		}

		name := ""
		if len(record) > nameIdx {
			name = strings.TrimSpace(record[nameIdx])
		}
		email := ""
		if len(record) > emailIdx {
			email = strings.TrimSpace(record[emailIdx])
		}

		if name == "" || email == "" {
			skipped++
			continue
		}

		// Check if user exists
		_, err = h.userRepo.FindByEmail(c.Request.Context(), email)
		if err == nil {
			skipped++
			skippedEmails = append(skippedEmails, email)
			continue // User exist
		}

		u := &domain.User{
			ID:     uuid.New(),
			RoleID: 3, // Peserta
			Name:   name,
			Email:  email,
		}

		if ageOk && len(record) > ageIdx {
			ageStr := strings.TrimSpace(record[ageIdx])
			if age, err := strconv.Atoi(ageStr); err == nil && age > 0 {
				u.Age = &age
			}
		}

		if posOk && len(record) > posIdx {
			pos := strings.TrimSpace(record[posIdx])
			if pos != "" {
				u.AppliedPosition = &pos
			}
		}

		if err := h.userRepo.Create(c.Request.Context(), u); err != nil {
			skipped++
			skippedEmails = append(skippedEmails, email)
		} else {
			imported++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Import selesai",
		"imported":       imported,
		"skipped":        skipped,
		"skipped_emails": skippedEmails,
	})
}

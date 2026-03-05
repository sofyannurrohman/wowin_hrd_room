package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/usecase"
	jwtpkg "hrd_room/backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC     *usecase.AuthUseCase
	jwtManager *jwtpkg.Manager
}

func NewAuthHandler(authUC *usecase.AuthUseCase, jwtManager *jwtpkg.Manager) *AuthHandler {
	return &AuthHandler{authUC: authUC, jwtManager: jwtManager}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUC.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": resp.Token,
		"user":  resp.User,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RoleID == 1 || req.RoleID == 2 {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required to create HR accounts"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := h.jwtManager.Validate(parts[1])
		if err != nil || claims.Role != domain.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "only super admin can create HR accounts"})
			return
		}
	}

	user, err := h.authUC.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userEmail, _ := c.Get("user_email")
	userRole, _ := c.Get("user_role")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   userEmail,
		"role":    userRole,
	})
}

// Apply handles public job applications (multipart/form-data)
func (h *AuthHandler) Apply(c *gin.Context) {
	// 1. Parse form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")
	ageStr := c.Request.FormValue("age")
	position := c.Request.FormValue("applied_position")
	expectedSalary := c.Request.FormValue("expected_salary")
	address := c.Request.FormValue("address")
	lastEdu := c.Request.FormValue("last_education")
	whatsapp := c.Request.FormValue("whatsapp_number")

	if name == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Email are required"})
		return
	}

	var age *int
	if ageStr != "" {
		parsedAge, err := strconv.Atoi(ageStr)
		if err == nil {
			age = &parsedAge
		}
	}

	var appliedPos *string
	if position != "" {
		appliedPos = &position
	}

	var expSalary *string
	if expectedSalary != "" {
		expSalary = &expectedSalary
	}

	var addr *string
	if address != "" {
		addr = &address
	}

	var edu *string
	if lastEdu != "" {
		edu = &lastEdu
	}

	var wa *string
	if whatsapp != "" {
		wa = &whatsapp
	}

	// 2. Handle CV Upload
	var cvURL *string
	file, header, err := c.Request.FormFile("cv")
	if err == nil {
		defer file.Close()

		// Create uploads folder if it doesn't exist
		uploadDir := "./uploads/cv"
		os.MkdirAll(uploadDir, os.ModePerm)

		// Generate unique filename
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%s-%s%s", strings.Split(email, "@")[0], hexEncodeRandom(4), ext)
		dstPath := filepath.Join(uploadDir, filename)

		out, err := os.Create(dstPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CV file"})
			return
		}
		defer out.Close()
		io.Copy(out, file)

		url := fmt.Sprintf("/uploads/cv/%s", filename)
		cvURL = &url
	}

	// 3. Generate Random Secure Password (users login via Magic Link anyway)
	password := hexEncodeRandom(8)

	req := usecase.RegisterRequest{
		Name:            name,
		Email:           email,
		Password:        password,
		RoleID:          3, // Peserta
		Age:             age,
		AppliedPosition: appliedPos,
		ExpectedSalary:  expSalary,
		Address:         addr,
		LastEducation:   edu,
		WhatsappNumber:  wa,
		CvURL:           cvURL,
	}

	user, err := h.authUC.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Application submitted successfully", "user_id": user.ID})
}

// helper
func hexEncodeRandom(bytes int) string {
	b := make([]byte, bytes)
	rand.Read(b)
	return hex.EncodeToString(b)
}

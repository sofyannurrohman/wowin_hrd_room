package handlers

import (
	"net/http"

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

	// Only allow Super Admin to register HR/Admin accounts
	userRole, _ := c.Get("user_role")
	if req.RoleID == 1 || req.RoleID == 2 {
		if userRole != domain.RoleSuperAdmin {
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

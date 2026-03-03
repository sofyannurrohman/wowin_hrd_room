package handlers

import (
	"net/http"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobPositionHandler struct {
	uc *usecase.JobPositionUseCase
}

func NewJobPositionHandler(uc *usecase.JobPositionUseCase) *JobPositionHandler {
	return &JobPositionHandler{uc: uc}
}

func (h *JobPositionHandler) ListActive(c *gin.Context) {
	positions, err := h.uc.GetActivePositions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"positions": positions})
}

func (h *JobPositionHandler) ListAll(c *gin.Context) {
	positions, err := h.uc.GetAllPositions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"positions": positions})
}

func (h *JobPositionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	p, err := h.uc.GetPosition(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "position not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"position": p})
}

func (h *JobPositionHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := &domain.JobPosition{
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if err := h.uc.CreatePosition(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "position created", "position": p})
}

func (h *JobPositionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := &domain.JobPosition{
		Name:     req.Name,
		IsActive: req.IsActive,
	}

	if err := h.uc.UpdatePosition(c.Request.Context(), id, p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "position updated"})
}

func (h *JobPositionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	if err := h.uc.DeletePosition(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "position deleted"})
}

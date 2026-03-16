package handlers

import (
	"net/http"

	"hrd_room/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardRepo *repository.DashboardRepository
}

func NewDashboardHandler(dashboardRepo *repository.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{dashboardRepo: dashboardRepo}
}

// GET /api/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashboardRepo.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GET /api/participants
func (h *DashboardHandler) ListParticipants(c *gin.Context) {
	participants, err := h.dashboardRepo.ListParticipants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, participants)
}

// GET /api/participants/:id
func (h *DashboardHandler) GetParticipantDetail(c *gin.Context) {
	participantID := c.Param("id")
	if participantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is missing"})
		return
	}

	detail, err := h.dashboardRepo.GetParticipantDetail(c.Request.Context(), participantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}

	c.JSON(http.StatusOK, detail)
}

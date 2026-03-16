package handlers

import (
	"net/http"
	"hrd_room/backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsUC *usecase.SettingsUseCase
}

func NewSettingsHandler(suc *usecase.SettingsUseCase) *SettingsHandler {
	return &SettingsHandler{settingsUC: suc}
}

func (h *SettingsHandler) Get(c *gin.Context) {
	key := c.Param("key")
	val, err := h.settingsUC.GetSettings(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, val)
}

func (h *SettingsHandler) Update(c *gin.Context) {
	key := c.Param("key")
	var val interface{}
	if err := c.ShouldBindJSON(&val); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handleError(err)})
		return
	}
	if err := h.settingsUC.UpdateSettings(c.Request.Context(), key, val); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handleError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}

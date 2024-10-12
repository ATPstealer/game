package controllers

import (
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSettings godoc
// @Summary Get General Game Settings
// @Description Get application settings
// @ID get-settings
// @Produce json
// @Success 200 {object} JSONResult
// @Router /api/v2/settings [get]
func GetSettings(c *gin.Context) {
	settings, err := models.GetSettings(db.M)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": settings})
}

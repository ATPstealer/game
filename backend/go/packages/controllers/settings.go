package controllers

import (
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": models.GetSettingsMap(db.DB)})
}

func GetSettingsMongo(c *gin.Context) {
	settings, err := models.GetSettingsMongo(db.M)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "code": 100001, "text": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": settings})
}

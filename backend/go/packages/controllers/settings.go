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

package controllers

import (
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSettings(c *gin.Context) {
	settings, err := models.GetSettings(db.M)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "code": 100001, "text": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": settings})
}

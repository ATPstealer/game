package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMyStorages(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	myStorages, err := models.GetMyStorages(db.M, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": "Can't get storages: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myStorages})
}

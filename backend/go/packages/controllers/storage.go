package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetMyStorages
//
//	@Summary	Return user's storages
//	@Tags		storage
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.Storage}
//	@Failure	401	{object}	JSONResult
//	@Failure	500	{object}	JSONResult
//	@Router		/storage/my [get]
func GetMyStorages(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	myStorages, err := models.GetMyStorages(db.M, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": "Can't get storages: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myStorages})
}

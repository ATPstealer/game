package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetResourceTypes
//
//	@Summary	Return all resource types from database
//	@Tags		resource
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.ResourceType}
//	@Failure	500	{object}	JSONResult
//	@Router		/resource/types [get]
func GetResourceTypes(c *gin.Context) {
	resourceTypes, err := models.GetAllResourceTypes(db.M)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": resourceTypes})
}

// GetMyResources
//
//	@Summary	Get user's resources
//	@Tags		resources
//	@Produce	json
//	@Param		x	query		int	false	"X coordinate"
//	@Param		y	query		int	false	"Y coordinate"
//	@Success	200	{object}	JSONResult{data=[]models.ResourceWithData}
//	@Failure	401	{object}	JSONResult
//	@Failure	500	{object}	JSONResult
//	@Router		/resource/my [get]
func GetMyResources(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	var xPointer, yPointer *int
	if c.Query("x") != "" {
		x, err := include.StrToInt(c, c.Query("x"))
		if err != nil {
			return
		}
		xPointer = &x
	}
	if c.Query("y") != "" {
		y, err := include.StrToInt(c, c.Query("y"))
		if err != nil {
			return
		}
		yPointer = &y
	}

	myResources, err := models.GetMyResources(db.M, userId, xPointer, yPointer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myResources})
}

// ResourceMove
//
//	@Summary	Initiates a resource movement
//	@Tags		logistics
//	@Accept		json
//	@Produce	json
//	@Param		logisticPayload	body		models.LogisticPayload	true	"Resource movement payload"
//	@Success	200				{object}	JSONResult
//	@Failure	401				{object}	JSONResult
//	@Failure	500				{object}	JSONResult
//	@Router		/resource/move [post]
func ResourceMove(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	var logisticPayload models.LogisticPayload
	if err := include.GetPayload(c, &logisticPayload); err != nil {
		return
	}

	if logisticPayload.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 21, "text": "Wrong amount"})
		return
	}

	err = models.StartLogisticJob(db.M, userId, logisticPayload)
	if err != nil {
		if strings.Contains(err.Error(), "not enough resources in this cell") {
			c.JSON(http.StatusOK, gin.H{"code": 22, "text": err.Error()})
		} else if strings.Contains(err.Error(), "there is not enough storage capacity in the destination sector") {
			c.JSON(http.StatusOK, gin.H{"code": 23, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough money") {
			c.JSON(http.StatusOK, gin.H{"code": 24, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough capacity in this hub") {
			c.JSON(http.StatusOK, gin.H{"code": 35, "text": err.Error()})
		} else if strings.Contains(err.Error(), "resource in different cell") {
			c.JSON(http.StatusOK, gin.H{"code": 37, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": -5, "values": logisticPayload})
}

// GetMyLogistics
//
//	@Summary	Get user's logistics tasks
//	@Tags		logistics
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.LogisticWithData}
//	@Failure	401	{object}	JSONResult
//	@Failure	500	{object}	JSONResult
//	@Router		/resource/my_logistics [get]
func GetMyLogistics(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	myLogistics, err := models.GetMyLogistics(db.M, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": "Can't get logistics: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myLogistics})
}

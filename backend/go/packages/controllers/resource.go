package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetResourceTypes(c *gin.Context) {
	resourceTypes, err := models.GetAllResourceTypes(db.DB)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get resource types: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": resourceTypes})
}

func GetMyResources(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
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

	myResources, err := models.GetMyResources(db.DB, userID, xPointer, yPointer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get resources: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myResources})
}

func ResourceMove(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	resourceTypeID, err := include.StrToUInt(c, c.Query("resource_type_id"))
	if err != nil {
		return
	}
	amount, err := include.StrTofloat64(c, c.Query("amount"))
	if err != nil {
		return
	}
	fromX, err := include.StrToInt(c, c.Query("from_x"))
	if err != nil {
		return
	}
	fromY, err := include.StrToInt(c, c.Query("from_y"))
	if err != nil {
		return
	}
	toX, err := include.StrToInt(c, c.Query("to_x"))
	if err != nil {
		return
	}
	toY, err := include.StrToInt(c, c.Query("to_y"))
	if err != nil {
		return
	}

	err = models.StartLogisticJob(db.DB, resourceTypeID, userID, amount, fromX, fromY, toX, toY)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't move resources: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Logistics company took you order for transfer " + fmt.Sprintf("%.1f", amount) +
		" from " + strconv.Itoa(fromX) + ":" + strconv.Itoa(fromY) + " to " + strconv.Itoa(toX) + ":" + strconv.Itoa(toY) + ". They didn't ask what was inside."})
}

func GetMyLogistics(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	myLogistics, err := models.GetMyLogistics(db.DB, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get logistics: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myLogistics})
}

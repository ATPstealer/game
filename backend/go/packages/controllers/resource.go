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

	var logisticPayload models.LogisticPayload
	if err := include.GetPayload(c, &logisticPayload); err != nil {
		return
	}

	if logisticPayload.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Wrong amount"})
		return
	}

	err = models.StartLogisticJob(db.DB, userID, logisticPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't move resources: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Logistics company took you order for transfer " + fmt.Sprintf("%.1f", logisticPayload.Amount) +
		" from " + strconv.Itoa(logisticPayload.FromX) + ":" + strconv.Itoa(logisticPayload.FromY) + " to " + strconv.Itoa(logisticPayload.ToX) + ":" +
		strconv.Itoa(logisticPayload.ToY) + ". They didn't ask what was inside."})
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

// mongo

func GetResourceTypesMongo(c *gin.Context) {
	resourceTypes, err := models.GetAllResourceTypesMongo(db.M)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get resource types: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": resourceTypes})
}

func GetMyResourcesMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
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

	myResources, err := models.GetMyResourcesMongo(db.M, userID, xPointer, yPointer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get resources: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myResources})
}

func ResourceMoveMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	var logisticPayload models.LogisticPayload
	if err := include.GetPayload(c, &logisticPayload); err != nil {
		return
	}

	if logisticPayload.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Wrong amount"})
		return
	}

	err = models.StartLogisticJobMongo(db.M, userID, logisticPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't move resources: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Logistics company took you order for transfer " + fmt.Sprintf("%.1f", logisticPayload.Amount) +
		" from " + strconv.Itoa(logisticPayload.FromX) + ":" + strconv.Itoa(logisticPayload.FromY) + " to " + strconv.Itoa(logisticPayload.ToX) + ":" +
		strconv.Itoa(logisticPayload.ToY) + ". They didn't ask what was inside."})
}

func GetMyLogisticsMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}
	myLogistics, err := models.GetMyLogisticsMongo(db.M, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get logistics: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myLogistics})
}

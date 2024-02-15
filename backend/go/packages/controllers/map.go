package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetMap(c *gin.Context) {
	mapCells := models.GetAllCells(db.DB)
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": mapCells})
}

func GetCellOwners(c *gin.Context) {
	x, err := include.StrToInt(c, c.Query("x"))
	if err != nil {
		return
	}
	y, err := include.StrToInt(c, c.Query("y"))
	if err != nil {
		return
	}
	cellOwners, err := models.GetCellOwners(db.DB, x, y)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get Cell Owners: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": cellOwners})
}

func GetAllLandLords(c *gin.Context) {
	cellOwners, err := models.GetAllLandLords(db.DB)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get Land Lords: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": cellOwners})
}

func GetMyLand(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	cellOwners, err := models.GetMyLands(db.DB, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get Land Lords: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": cellOwners})
}

func BuyLand(c *gin.Context) {
	var buyLandPayload models.BuyLandPayload
	var err error
	if err = include.GetPayload(c, &buyLandPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	price, err := models.BuyLand(db.DB, userID, buyLandPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't buy land: " + err.Error()})
		log.Println("Can't buy land: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You bought %d ares in %dx%d cell by %.2f$",
		buyLandPayload.Square, buyLandPayload.X, buyLandPayload.Y, price)})
}

func GetCellOwnersMongo(c *gin.Context) {
	x, err := include.StrToInt(c, c.Query("x"))
	if err != nil {
		return
	}
	y, err := include.StrToInt(c, c.Query("y"))
	if err != nil {
		return
	}
	cellOwners, err := models.GetCellOwnersMongo(db.M, x, y)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 15, "text": "Can't get Cell Owners: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": cellOwners})
}

// mongo

func BuyLandMongo(c *gin.Context) {
	var buyLandPayload models.BuyLandPayload
	var err error
	if err = include.GetPayload(c, &buyLandPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	price, err := models.BuyLandMongo(db.M, userID, buyLandPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 16, "text": "Can't buy land: " + err.Error()})
		log.Println("Can't buy land: " + err.Error())
		return
	}

	values := map[string]interface{}{
		"square": buyLandPayload.Square,
		"x":      buyLandPayload.X,
		"y":      buyLandPayload.Y,
		"price":  price,
	}

	valuesJson, err := json.Marshal(values)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 17, "text": "Can't make JSON: " + err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "code": -3, "values": string(valuesJson),
		"text": fmt.Sprintf("You bought %d ares in %dx%d cell by %.2f$",
			buyLandPayload.Square, buyLandPayload.X, buyLandPayload.Y, price)})
}

func GetMapMongo(c *gin.Context) {
	mapCells, err := models.GetAllCellsMongo(db.M)
	if err != nil {
		log.Println("Can't get map cells: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 18, "text": "Can't get map cells: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": mapCells})
}

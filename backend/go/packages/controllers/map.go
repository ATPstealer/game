package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func GetCellOwners(c *gin.Context) {
	x, err := include.StrToInt(c, c.Query("x"))
	if err != nil {
		return
	}
	y, err := include.StrToInt(c, c.Query("y"))
	if err != nil {
		return
	}
	cellOwners, err := models.GetCellOwners(db.M, x, y)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cellOwners})
}

func BuyLand(c *gin.Context) {
	var buyLandPayload models.BuyLandPayload

	var err error
	if err = include.GetPayload(c, &buyLandPayload); err != nil {
		return
	}

	if buyLandPayload.Square <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 26})
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	price, err := models.BuyLand(db.M, userId, buyLandPayload)
	if err != nil {
		if strings.Contains(err.Error(), "square should be greater than 0") {
			c.JSON(http.StatusOK, gin.H{"code": 26, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough money") {
			c.JSON(http.StatusOK, gin.H{"code": 24, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough land in this cell") {
			c.JSON(http.StatusOK, gin.H{"code": 27, "text": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
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
		log.Println("Can't make JSON: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 100011, "text": "Can't make JSON: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -3, "values": string(valuesJson)})
}

func GetMap(c *gin.Context) {
	mapCells, err := models.GetAllCells(db.M)
	if err != nil {
		log.Println("Can't get map cells: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": mapCells})
}

func GetAllLandLords(c *gin.Context) {
	cellOwners, err := models.GetAllLandLords(db.M)
	if err != nil {
		log.Println("Can't get Land Lords: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cellOwners})
}

func GetMyLand(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	myLands, err := models.GetMyLands(db.M, userId)
	if err != nil {
		log.Println("Can't get Land Lords: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myLands})
}

package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStoreGoodsMongo(c *gin.Context) {
	buildingID, err := include.StrToPrimObjId(c, c.Query("building_id")) // TODO: Сделать buildingId
	if err != nil {
		return
	}

	storePrices, err := models.GetStoreGoodsMongo(db.M, buildingID)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get store prices: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": storePrices})
}

func SetStoreGoodsMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	var storeGoodsPayload models.StoreGoodsPayloadMongo
	if err := include.GetPayload(c, &storeGoodsPayload); err != nil {
		return
	}

	err = models.SetStoreGoodsMongo(db.M, userID, storeGoodsPayload)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't set price: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok"})
}

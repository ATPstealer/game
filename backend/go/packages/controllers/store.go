package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStoreGoods(c *gin.Context) {
	buildingID, err := include.StrToPrimObjId(c, c.Query("buildingId"))
	if err != nil {
		return
	}

	storePrices, err := models.GetStoreGoods(db.M, buildingID)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get store prices: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": storePrices})
}

func SetStoreGoods(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	var storeGoodsPayload models.StoreGoodsPayload
	if err := include.GetPayload(c, &storeGoodsPayload); err != nil {
		return
	}

	err = models.SetStoreGoods(db.M, userID, storeGoodsPayload)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't set price: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok"})
}

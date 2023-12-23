package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStoreGoods(c *gin.Context) {
	buildingID, err := include.StrToUInt(c, c.Query("building_id"))
	if err != nil {
		return
	}

	storePrices, err := models.GetStoreGoods(db.DB, buildingID)

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

	buildingID, err := include.StrToUInt(c, c.Query("building_id"))
	if err != nil {
		return
	}
	resourceTypeID, err := include.StrToUInt(c, c.Query("resource_type_id"))
	if err != nil {
		return
	}
	price, err := include.StrToFloat32(c, c.Query("price"))
	if err != nil {
		return
	}

	err = models.SetStoreGoods(db.DB, userID, buildingID, resourceTypeID, price)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't set price: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok"})
}

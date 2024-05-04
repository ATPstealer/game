package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SetStoreGoods(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	var storeGoodsPayload models.StoreGoodsPayload
	if err := include.GetPayload(c, &storeGoodsPayload); err != nil {
		return
	}

	err = models.SetStoreGoods(db.M, userId, storeGoodsPayload)

	if err != nil {
		if strings.Contains(err.Error(), "this building doesn't belong to you") {
			c.JSON(http.StatusOK, gin.H{"code": 18, "text": err.Error()})
		} else if strings.Contains(err.Error(), "this is not a store") {
			c.JSON(http.StatusOK, gin.H{"code": 19, "text": err.Error()})
		} else if strings.Contains(err.Error(), "can't sell here") {
			c.JSON(http.StatusOK, gin.H{"code": 20, "text": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0})
}

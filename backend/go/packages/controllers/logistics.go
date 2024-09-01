package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetLogisticsCapacity(c *gin.Context) {
	var findLogisticsParams models.FindLogisticsParams

	if c.Query("x") != "" {
		x, err := include.StrToInt(c, c.Query("x"))
		if err != nil {
			return
		}
		findLogisticsParams.X = &x
	}
	if c.Query("y") != "" {
		y, err := include.StrToInt(c, c.Query("y"))
		if err != nil {
			return
		}
		findLogisticsParams.Y = &y
	}
	if c.Query("minCapacity") != "" {
		minCapacity, err := include.StrToFloat64(c, c.Query("minCapacity"))
		if err != nil {
			return
		}
		findLogisticsParams.MinCapacity = &minCapacity
	}

	logisticsCapacity, err := models.GetLogisticsCapacity(db.M, findLogisticsParams)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": logisticsCapacity})
}

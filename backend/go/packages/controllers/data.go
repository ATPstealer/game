package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetEvolutionPrices(c *gin.Context) {
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

	evolutionPrices, err := models.GetEvolutionPrices(db.M, xPointer, yPointer)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get evolution prices: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": evolutionPrices})
}

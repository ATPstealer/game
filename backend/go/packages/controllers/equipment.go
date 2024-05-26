package controllers

import (
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetEquipmentTypes(c *gin.Context) {
	equipmentTypes, err := models.GetAllEquipmentTypes(db.M)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": equipmentTypes})
}

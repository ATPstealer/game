package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetEquipmentTypes
//
//	@Summary	Get all equipment types
//	@Tags		equipment
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.EquipmentType}
//	@Failure	500	{object}	JSONResult
//	@Router		/equipment/types [get]
func GetEquipmentTypes(c *gin.Context) {
	equipmentTypes, err := models.GetAllEquipmentTypes(db.M)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": equipmentTypes})
}

// GetMyEquipment
//
//	@Summary	Return user's equipment
//	@Tags		equipment
//	@Accept		json
//	@Produce	json
//	@Param		x	query		int	false	"x-coordinate of the equipment location"
//	@Param		y	query		int	false	"y-coordinate of the equipment location"
//	@Success	200	{object}	JSONResult{data=[]models.ResourceAsEquipment}
//	@Failure	401	{object}	JSONResult
//	@Failure	500	{object}	JSONResult
//	@Router		/equipment/my [get]
func GetMyEquipment(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
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

	equipment, err := models.GetMyEquipment(db.M, userId, xPointer, yPointer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": equipment})
}

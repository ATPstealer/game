package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// GetCellOwners
//
//	@Summary	Get the landlords in cell
//	@Tags		map
//	@Accept		json
//	@Produce	json
//	@Param		x	query		int	true	"X coordinate"
//	@Param		y	query		int	true	"Y coordinate"
//	@Success	200	{object}	JSONResult{data=[]models.LandLord}
//	@Failure	500	{object}	JSONResult
//	@Router		/map/cell_owners [get]
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cellOwners})
}

// BuyLand
//
//	@Summary	Buy land in cell
//	@Tags		map
//	@Accept		json
//	@Produce	json
//	@Param		buyLandPayload	body		models.BuyLandPayload	true	"Land purchase payload"
//	@Success	200				{object}	JSONResult{values=models.BuyLandPayload}
//	@Failure	500				{object}	JSONResult
//	@Router		/map/buy_land [post]
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
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	values := map[string]interface{}{
		"square": buyLandPayload.Square,
		"x":      buyLandPayload.X,
		"y":      buyLandPayload.Y,
		"price":  price,
	}

	c.JSON(http.StatusOK, gin.H{"code": -3, "values": values})
}

// GetMap
//
//	@Summary		Return map cells
//	@Tags			map
//	@Description	Returns the list of all map cells
//	@Produce		json
//	@Success		200	{object}	JSONResult{data=[]models.Cell}
//	@Failure		500	{object}	JSONResult
//	@Router			/map [get]
func GetMap(c *gin.Context) {
	mapCells, err := models.GetAllCells(db.M)
	if err != nil {
		log.Println("Can't get map cells: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": mapCells})
}

// GetAllLandLords
//
//	@Summary	Return all landowners
//	@Tags		map
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.LandLord}
//	@Failure	500	{object}	JSONResult
//	@Router		/map/all_land_lords [get]
func GetAllLandLords(c *gin.Context) {
	cellOwners, err := models.GetAllLandLords(db.M)
	if err != nil {
		log.Println("Can't get Land Lords: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cellOwners})
}

// GetMyLand
//
//	@Summary	Return user's lands
//	@Tags		map
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.LandLord}
//	@Failure	500	{object}	JSONResult
//	@Router		/map/my [get]
func GetMyLand(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	myLands, err := models.GetMyLands(db.M, userId)
	if err != nil {
		log.Println("Can't get Land Lords: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myLands})
}

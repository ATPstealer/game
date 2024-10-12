package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
)

// GetBlueprints
//
//	@Summary		Get blueprints
//	@Description	Fetches a list of blueprints. If an 'id' query parameter is provided, fetches the blueprint with the specified ID.
//	@Tags			blueprints
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	false	"Blueprint ID"
//	@Success		200	{object}	JSONResult{data=[]models.Blueprint}
//	@Failure		500	{object}	JSONResult
//	@Router			/building/blueprints [get]
func GetBlueprints(c *gin.Context) {
	var blueprintId uint
	var err error
	if c.Query("id") != "" {
		blueprintId, err = include.StrToUInt(c, c.Query("id"))
		if err != nil {
			return
		}
	}
	blueprints, err := models.GetBlueprints(db.M, blueprintId)
	if err != nil {
		log.Printf("Can't get blueprints: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": blueprints})
}

// GetBuildingsTypes
//
//	@Summary	Get all building types
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.BuildingType}
//	@Failure	500	{object}	JSONResult
//	@Router		/building/types [get]
func GetBuildingsTypes(c *gin.Context) {
	buildingTypes, err := models.GetAllBuildingTypes(db.M)
	if err != nil {
		log.Println("Can't get building types: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": buildingTypes})
}

// ConstructBuilding
//
//	@Summary	Construct a new building
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Param		constructBuildingPayload	body		models.ConstructBuildingPayload	true	"Building construction payload"
//	@Success	200							{object}	JSONResult{values=models.ConstructBuildingPayload}
//	@Failure	401							{object}	JSONResult
//	@Failure	500							{object}	JSONResult
//	@Router		/building/construct [post]
func ConstructBuilding(c *gin.Context) {
	var constructBuildingPayload models.ConstructBuildingPayload
	if err := include.GetPayload(c, &constructBuildingPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	if constructBuildingPayload.Square <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 26})
		return
	}

	err = models.ConstructBuilding(db.M, userId, constructBuildingPayload)
	if err != nil {
		if strings.Contains(err.Error(), "not enough land in this cell") {
			c.JSON(http.StatusOK, gin.H{"code": 27, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough money") {
			c.JSON(http.StatusOK, gin.H{"code": 24, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -9, "values": constructBuildingPayload})
}

// GetBuildings
//
//	@Summary	Fetch the list of buildings
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Param		findBuildingsParams	body		models.FindBuildingParams	false	"Parameters to filter and sort buildings"
//	@Success	200					{object}	JSONResult{values=[]models.BuildingWithData}
//	@Failure	500					{object}	JSONResult
//	@Router		/building/get [post]
func GetBuildings(c *gin.Context) {
	var findBuildingsParams models.FindBuildingParams
	if err := include.GetPayload(c, &findBuildingsParams); err != nil {
		return
	}

	if findBuildingsParams.NickName != nil {
		User, err := models.GetUserByNickName(db.M, *findBuildingsParams.NickName)
		if err != nil {
			log.Printf("Can't get user by nickname: " + err.Error())
			c.JSON(http.StatusOK, gin.H{"code": 13, "text": "Can't get user by nickname: " + err.Error()})
			return
		}
		findBuildingsParams.UserId = &User.Id
	}

	buildings, err := models.GetBuildings(db.M, findBuildingsParams)

	if err != nil {
		log.Printf("Can't get buildings: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": buildings})
}

// GetMyBuildings
//
//	@Summary		Fetch the user's buildings
//	@Tags			buildings
//	@Description	Optionally filter by building ID.
//	@Accept			json
//	@Produce		json
//	@Param			_id	query		string	false	"Building ID to filter by"
//	@Success		200	{object}	JSONResult{data=[]models.BuildingWithData}
//	@Failure		401	{object}	JSONResult
//	@Failure		500	{object}	JSONResult
//	@Router			/building/my [get]
func GetMyBuildings(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	var buildingId primitive.ObjectID
	if c.Query("_id") != "" {
		buildingId, err = include.StrToPrimObjId(c, c.Query("_id"))
		if err != nil {
			return
		}
	}
	myBuildings, err := models.GetMyBuildings(db.M, userId, buildingId)
	if err != nil {
		log.Printf("Can't get my buildings: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myBuildings})
}

// StartWork
//
//	@Summary	Start work in the building
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Param		startWorkPayload	body		models.StartProductionPayload	true	"Start production payload"
//	@Success	200					{object}	JSONResult
//	@Failure	401					{object}	JSONResult
//	@Failure	500					{object}	JSONResult
//	@Router		/building/start_work [post]
func StartWork(c *gin.Context) {
	var startWorkPayload models.StartProductionPayload
	var err error
	if err = include.GetPayload(c, &startWorkPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.StartProduction(db.M, userId, startWorkPayload)
	if err != nil {
		if strings.Contains(err.Error(), "building busy") {
			c.JSON(http.StatusOK, gin.H{"code": 28, "text": err.Error()})
		} else if strings.Contains(err.Error(), "this building don't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else if strings.Contains(err.Error(), "can't product it here") {
			c.JSON(http.StatusOK, gin.H{"code": 30, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -10})
}

// StopWork
//
//	@Summary	Stops any work in building. Later he should stop only the works available for stopping.
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Param		startWorkPayload	body		models.StartProductionPayload	true	"Production stop payload"
//	@Success	200					{object}	JSONResult
//	@Failure	401					{object}	JSONResult
//	@Failure	500					{object}	JSONResult
//	@Router		/building/stop_work [post]
func StopWork(c *gin.Context) {
	var startWorkPayload models.StartProductionPayload
	var err error
	if err = include.GetPayload(c, &startWorkPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.StopProduction(db.M, userId, startWorkPayload)

	if err != nil {
		if strings.Contains(err.Error(), "this building don't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -13})
}

// SetHiring
//
//	@Summary	Set hiring details for a building
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Param		hiringPayload	body		models.HiringPayload	true	"Details of hiring"
//	@Success	200				{object}	JSONResult
//	@Failure	401				{object}	JSONResult
//	@Failure	500				{object}	JSONResult
//	@Router		/building/hiring [post]
func SetHiring(c *gin.Context) {
	var hiringPayload models.HiringPayload
	if err := include.GetPayload(c, &hiringPayload); err != nil {
		return
	}
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	if hiringPayload.HiringNeeds < 0 || hiringPayload.Salary < 0 {
		c.JSON(http.StatusOK, gin.H{"code": "21"})
		return
	}

	err = models.SetHiring(db.M, userId, hiringPayload)
	if err != nil {
		if strings.Contains(err.Error(), "his building doesn't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else if strings.Contains(err.Error(), "hiring needs more that maximum") {
			c.JSON(http.StatusOK, gin.H{"code": 31, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -11})
}

// DestroyBuilding
//
//	@Summary	Destroy an existing building
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Param		_id	query		string	true	"Building ID"
//	@Success	200	{object}	JSONResult
//	@Failure	401	{object}	JSONResult
//	@Failure	500	{object}	JSONResult
//	@Router		/building/destroy [delete]
func DestroyBuilding(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	buildingID, err := include.StrToPrimObjId(c, c.Query("_id"))
	if err != nil {
		return
	}
	err = models.DestroyBuilding(db.M, userId, buildingID)
	if err != nil {
		if strings.Contains(err.Error(), "for attempting to destroy someone else's building, inevitable punishment awaits you") {
			c.JSON(http.StatusOK, gin.H{"code": 32, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -12})
}

// InstallEquipment
//
//	@Summary	Install equipment in a building
//	@Tags		buildings
//	@Accept		json
//	@Produce	json
//	@Param		installEquipmentPayload	body		models.InstallEquipmentPayload	true	"Equipment installation payload"
//	@Success	200						{object}	JSONResult
//	@Failure	401						{object}	JSONResult
//	@Failure	500						{object}	JSONResult
//	@Router		/building/install_equipment [post]
func InstallEquipment(c *gin.Context) {
	var installEquipmentPayload models.InstallEquipmentPayload
	if err := include.GetPayload(c, &installEquipmentPayload); err != nil {
		return
	}
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	if err := models.InstallEquipment(db.M, userId, installEquipmentPayload); err != nil {
		if strings.Contains(err.Error(), "this building don't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough resources in this cell") {
			c.JSON(http.StatusOK, gin.H{"code": 22, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough equipment here") {
			c.JSON(http.StatusOK, gin.H{"code": 33, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough space") {
			c.JSON(http.StatusOK, gin.H{"code": 34, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": -14})
}

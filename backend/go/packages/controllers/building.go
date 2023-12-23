package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetBuildingsTypes(c *gin.Context) {
	buildingTypes, err := models.GetAllBuildingTypes(db.DB)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get building types: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": buildingTypes})
}

func GetBuildings(c *gin.Context) {
	var findBuildingsParams models.FindBuildingParams
	if err := include.GetPayload(c, &findBuildingsParams); err != nil {
		return
	}

	if findBuildingsParams.NickName != nil {
		User, err := models.GetUserByNickName(db.DB, *findBuildingsParams.NickName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get user: " + err.Error()})
			return
		}
		findBuildingsParams.UserID = &User.ID
	}

	buildings, err := models.GetBuildings(db.DB, findBuildingsParams)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get buildings: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": buildings})
}

func GetBlueprints(c *gin.Context) {
	var blueprintID uint
	var err error
	if c.Query("id") != "" {
		blueprintID, err = include.StrToUInt(c, c.Query("id"))
		if err != nil {
			return
		}
	}
	blueprints, err := models.GetBlueprints(db.DB, blueprintID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Can't get blueprints: " + err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": blueprints})
}

func CreateBuilding(c *gin.Context) {
	var constructBuildingPayload models.ConstructBuildingPayload
	if err := include.GetPayload(c, &constructBuildingPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	if constructBuildingPayload.Square <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "In this room, people will suffer from lack of air."})
		return
	}

	err = models.ConstructBuilding(db.DB, userID, constructBuildingPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't create building: " + err.Error()})
		log.Println("Can't create building: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You start to construct building")})
}

func GetMyBuildings(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	var buildingID uint
	if c.Query("id") != "" {
		buildingID, err = include.StrToUInt(c, c.Query("id"))
		if err != nil {
			return
		}
	}
	myBuildings, err := models.GetMyBuildings(db.DB, userID, buildingID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get buildings: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myBuildings})
}

func StartWork(c *gin.Context) {
	var startWorkPayload models.StartWorkPayload
	var err error
	if err = include.GetPayload(c, &startWorkPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	err = models.StartWork(db.DB, userID, startWorkPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't start job: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Job was stared"})
}

func DestroyBuilding(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	log.Println(c.Query("id"))

	buildingID, err := include.StrToUInt(c, c.Query("id"))
	if err != nil {
		return
	}
	err = models.DestroyBuilding(db.DB, userID, buildingID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't destroy building: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Building has been destroyed :("})

}

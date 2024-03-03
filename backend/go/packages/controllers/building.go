package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func GetBlueprintsMongo(c *gin.Context) {
	var blueprintID uint
	var err error
	if c.Query("id") != "" {
		blueprintID, err = include.StrToUInt(c, c.Query("id"))
		if err != nil {
			return
		}
	}
	blueprints, err := models.GetBlueprintsMongo(db.M, blueprintID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 11, "text": "Can't get blueprints: " + err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": blueprints})
}

func GetBuildingsTypesMongo(c *gin.Context) {
	buildingTypes, err := models.GetAllBuildingTypesMongo(db.M)
	if err != nil {
		log.Println("Can't get building types: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 1, "text": "Can't get building types: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": buildingTypes})
}

func ConstructBuildingMongo(c *gin.Context) {
	var constructBuildingPayload models.ConstructBuildingPayload
	if err := include.GetPayload(c, &constructBuildingPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	if constructBuildingPayload.Square <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "In this room, people will suffer from lack of air."})
		return
	}

	err = models.ConstructBuildingMongo(db.M, userID, constructBuildingPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't create building: " + err.Error()})
		log.Println("Can't create building: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You start to construct building")})
}

func GetBuildingsMongo(c *gin.Context) {
	var findBuildingsParamsMongo models.FindBuildingParamsMongo
	if err := include.GetPayload(c, &findBuildingsParamsMongo); err != nil {
		return
	}

	if findBuildingsParamsMongo.NickName != nil {
		User, err := models.GetUserByNickNameMongo(db.M, *findBuildingsParamsMongo.NickName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 13, "text": "Can't get user: " + err.Error()})
			return
		}
		findBuildingsParamsMongo.UserID = &User.ID
	}

	buildings, err := models.GetBuildingsMongo(db.M, findBuildingsParamsMongo)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 19, "text": "Can't get buildings: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": buildings})
}

func GetMyBuildingsMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}
	var buildingID primitive.ObjectID
	if c.Query("_id") != "" {
		buildingID, err = include.StrToPrimObjId(c, c.Query("_id"))
		if err != nil {
			return
		}
	}
	myBuildings, err := models.GetMyBuildingsMongo(db.M, userID, buildingID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 1, "text": "Can't get buildings: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": myBuildings})
}

func StartWorkMongo(c *gin.Context) {
	var startWorkPayload models.StartWorkPayloadMongo
	var err error
	if err = include.GetPayload(c, &startWorkPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	err = models.StartWorkMongo(db.M, userID, startWorkPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't start job: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Job was stared"})
}

func SetHiringMongo(c *gin.Context) {
	var hiringPayload models.HiringPayloadMongo
	if err := include.GetPayload(c, &hiringPayload); err != nil {
		return
	}
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}
	if hiringPayload.HiringNeeds < 0 || hiringPayload.Salary < 0 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't set < 0"})
		return
	}

	err = models.SetHiringMongo(db.M, userID, hiringPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't change hiring details: " + err.Error()})
		log.Println("Can't change hiring details: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("Hiring details changed")})
}

func DestroyBuildingMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	buildingID, err := include.StrToPrimObjId(c, c.Query("_id"))
	if err != nil {
		return
	}
	err = models.DestroyBuildingMongo(db.M, userID, buildingID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't destroy building: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "Building has been destroyed :("})
}

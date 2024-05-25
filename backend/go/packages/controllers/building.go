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
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": blueprints})
}

func GetBuildingsTypes(c *gin.Context) {
	buildingTypes, err := models.GetAllBuildingTypes(db.M)
	if err != nil {
		log.Println("Can't get building types: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": buildingTypes})
}

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
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -9, "values": constructBuildingPayload})
}

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
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": buildings})
}

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
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myBuildings})
}

func StartWork(c *gin.Context) {
	var startWorkPayload models.StartWorkPayload
	var err error
	if err = include.GetPayload(c, &startWorkPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.StartWork(db.M, userId, startWorkPayload)
	if err != nil {
		if strings.Contains(err.Error(), "building busy") {
			c.JSON(http.StatusOK, gin.H{"code": 28, "text": err.Error()})
		} else if strings.Contains(err.Error(), "this building don't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else if strings.Contains(err.Error(), "can't product it here") {
			c.JSON(http.StatusOK, gin.H{"code": 30, "text": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -10})
}

func StopWork(c *gin.Context) {
	var startWorkPayload models.StartWorkPayload
	var err error
	if err = include.GetPayload(c, &startWorkPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.StopWork(db.M, userId, startWorkPayload)

	if err != nil {
		if strings.Contains(err.Error(), "this building don't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -13})
}

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
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -11})
}

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
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -12})
}

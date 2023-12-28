package models

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type BuildingType struct {
	gorm.Model                     // This includes some common fields like ID, CreatedAt, UpdatedAt, and DeletedAt.
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	Cost             float64       `json:"cost"`
	Requirements     string        `json:"requirements"` // TODO: Resource object
	BuildTime        time.Duration `json:"buildTime"`
	BuildingGroup    string        `json:"buildingGroup"`
	BuildingSubGroup string        `json:"buildingSubGroup"`
	Capacity         float64       `json:"capacity"`
}

type BuildingTypeResult struct {
	ID               uint          `json:"id"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	Cost             float64       `json:"cost"`
	Requirements     string        `json:"requirements"`
	BuildTime        time.Duration `json:"buildTime"`
	BuildingGroup    string        `json:"buildingGroup"`
	BuildingSubGroup string        `json:"buildingSubGroup"`
	Capacity         float64       `json:"capacity"`
}

func GetBuildingTypeByID(db *gorm.DB, typeID uint) (BuildingTypeResult, error) {
	var buildingType BuildingTypeResult
	res := db.Model(&BuildingType{}).Where("id = ?", typeID).First(&buildingType)
	if res.Error != nil {
		log.Println("Can't get Building Type: " + res.Error.Error())
	}
	return buildingType, res.Error
}

func GetAllBuildingTypes(db *gorm.DB) ([]BuildingTypeResult, error) {
	var buildingTypes []BuildingTypeResult
	res := db.Model(&BuildingType{}).Find(&buildingTypes).Select("title")
	if res.Error != nil {
		log.Println("Can't get Building Type: " + res.Error.Error())
	}
	return buildingTypes, res.Error
}

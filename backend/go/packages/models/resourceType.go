package models

import (
	"gorm.io/gorm"
	"log"
)

type ResourceType struct {
	gorm.Model
	Name       string  `json:"name"`
	Volume     float32 `json:"volume"` // m3
	Weight     float32 `json:"weight"` // kg
	Demand     float32 `json:"demand"`
	StoreGroup string  `json:"storeGroup"`
}

type ResourceTypeResult struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Volume     float32 `json:"volume"`
	Weight     float32 `json:"weight"`
	Demand     float32 `json:"demand"`
	StoreGroup string  `json:"storeGroup"`
}

func GetAllResourceTypes(db *gorm.DB) ([]ResourceTypeResult, error) {
	var resourceTypes []ResourceTypeResult
	res := db.Model(&ResourceType{}).Find(&resourceTypes)
	if res.Error != nil {
		log.Println("Can't get Resource Types: " + res.Error.Error())
	}
	return resourceTypes, res.Error
}

func GetResourceTypesByID(db *gorm.DB, typeID uint) (ResourceTypeResult, error) {
	var resourceType ResourceTypeResult
	res := db.Model(&ResourceType{}).Where("id = ?", typeID).First(&resourceType)
	if res.Error != nil {
		log.Println("Can't get Resource Type: " + res.Error.Error())
	}
	return resourceType, res.Error
}

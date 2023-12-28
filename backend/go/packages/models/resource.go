package models

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type Resource struct {
	gorm.Model
	ResourceTypeID uint    `json:"resourceTypeId"`
	UserID         uint    `json:"userId"`
	Amount         float64 `json:"amount"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
}

type ResourceResult struct {
	ID             uint    `json:"id"`
	UserID         uint    `json:"userId"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	Amount         float64 `json:"amount"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	Name           string  `json:"name"`
	Volume         float64 `json:"volume"`
	Weight         float64 `json:"weight"`
}

func GetAllResources(db *gorm.DB) ([]ResourceResult, error) {
	var allResources []ResourceResult
	res := db.Model(&Resource{}).
		Select("resources.id", "user_id", "resource_type_id", "amount", "x", "y", "name", "volume", "weight").
		Joins("left join resource_types on resources.resource_type_id = resource_types.id").
		Scan(&allResources)

	if res.Error != nil {
		log.Println("Can't get resources: " + res.Error.Error())
	}
	return allResources, res.Error
}

func GetMyResources(db *gorm.DB, userID uint) ([]ResourceResult, error) {
	var resources []ResourceResult
	res := db.Model(&Resource{}).Where("user_id", userID).
		Select("resources.id", "resource_type_id", "amount", "x", "y", "name", "volume", "weight").
		Joins("left join resource_types on resources.resource_type_id = resource_types.id").
		Scan(&resources)

	if res.Error != nil {
		log.Println("Can't get resources: " + res.Error.Error())
	}
	return resources, res.Error
}

func GetMyResourceInCell(db *gorm.DB, resourceTypeID uint, userID uint, x int, y int) (ResourceResult, error) {
	var resource ResourceResult
	res := db.Model(&Resource{}).Where("user_id = ? AND resource_type_id = ? AND X = ? AND Y = ?",
		userID, resourceTypeID, x, y).
		Select("resources.id", "resource_type_id", "amount", "x", "y", "name", "volume", "weight").
		Joins("left join resource_types on resources.resource_type_id = resource_types.id").
		First(&resource)

	if res.Error != nil {
		log.Println("Can't get resources: " + res.Error.Error())
	}
	return resource, res.Error
}

func AddResource(db *gorm.DB, resourceTypeID uint, userID uint, x int, y int, amount float64) error {
	newResource := Resource{
		ResourceTypeID: resourceTypeID,
		UserID:         userID,
		X:              x,
		Y:              y,
	}
	result := db.Model(&Resource{}).Where("resource_type_id =? AND user_id = ? AND x = ? AND y = ?", resourceTypeID, userID, x, y).
		First(&newResource)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newResource.Amount = amount
			db.Create(&newResource)
		} else {
			return result.Error
		}
	} else {
		newResource.Amount += amount
		db.Save(&newResource)
	}
	return nil
}

func CheckEnoughResources(db *gorm.DB, resourceTypeID uint, userID uint, x int, y int, amount float64) bool {
	var resource Resource
	result := db.Model(&Resource{}).Where("resource_type_id =? AND user_id = ? AND x = ? AND y = ?", resourceTypeID, userID, x, y).
		First(&resource)
	if result.Error != nil {
		return false
	}
	return resource.Amount >= amount
}

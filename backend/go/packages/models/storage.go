package models

import (
	"gorm.io/gorm"
	"log"
)

type Storage struct {
	gorm.Model
	UserID         uint    `json:"userId"`
	VolumeOccupied float32 `json:"volumeOccupied"`
	VolumeMax      float32 `json:"volumeMax"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
}

type StorageResult struct {
	gorm.Model
	UserID         uint    `json:"userId"`
	VolumeOccupied float32 `json:"volumeOccupied"`
	VolumeMax      float32 `json:"volumeMax"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
}

func GetMyStorages(db *gorm.DB, userID uint) ([]StorageResult, error) {
	var myStorages []StorageResult
	res := db.Model(&Storage{}).Where("user_id", userID).Scan(&myStorages)
	if res.Error != nil {
		log.Println("Can't get buildings: " + res.Error.Error())
	}
	return myStorages, res.Error
}

func GetAllStorages(db *gorm.DB) ([]Storage, error) {
	var allStorages []Storage
	res := db.Model(&Storage{}).Scan(&allStorages)
	if res.Error != nil {
		log.Println("Can't get buildings: " + res.Error.Error())
	}
	return allStorages, res.Error
}

func CheckEnoughStorage(db *gorm.DB, userID uint, x int, y int, addVolume float32) bool {
	var storage Storage
	res := db.Model(&Storage{}).Where("user_id = ? AND x = ? AND y = ?", userID, x, y).First(&storage)
	if res.Error != nil {
		log.Println("Can't get storages: " + res.Error.Error())
		return false
	}
	return storage.VolumeMax >= storage.VolumeOccupied+addVolume+GetDestinationVolume(db, userID, x, y)
}

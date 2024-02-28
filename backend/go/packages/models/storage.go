package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
)

type Storage struct {
	gorm.Model
	UserID         uint    `json:"userId"`
	VolumeOccupied float64 `json:"volumeOccupied"`
	VolumeMax      float64 `json:"volumeMax"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
}

type StorageResult struct {
	gorm.Model
	UserID         uint    `json:"userId"`
	VolumeOccupied float64 `json:"volumeOccupied"`
	VolumeMax      float64 `json:"volumeMax"`
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

func CheckEnoughStorage(db *gorm.DB, userID uint, x int, y int, addVolume float64) bool {
	var storage Storage
	res := db.Model(&Storage{}).Where("user_id = ? AND x = ? AND y = ?", userID, x, y).First(&storage)
	if res.Error != nil {
		log.Println("Can't get storages: " + res.Error.Error())
		return false
	}
	return storage.VolumeMax >= storage.VolumeOccupied+addVolume+GetDestinationVolume(db, userID, x, y)
}

// mongo

type StorageMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	VolumeOccupied float64            `json:"volumeOccupied" bson:"volumeOccupied"`
	VolumeMax      float64            `json:"volumeMax" bson:"volumeMax"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
}

func GetAllStoragesMongo(m *mongo.Database) ([]StorageMongo, error) {
	var allStorages []StorageMongo
	cursor, err := m.Collection("storages").Find(context.TODO(), bson.M{})
	if err != nil {
		return allStorages, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &allStorages)
	return allStorages, err
}

func GetMyStoragesMongo(m *mongo.Database, userID primitive.ObjectID) ([]StorageMongo, error) {
	var allStorages []StorageMongo
	cursor, err := m.Collection("storages").Find(context.TODO(), bson.M{"userId": userID})
	if err != nil {
		return allStorages, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &allStorages)
	return allStorages, err
}

func CheckEnoughStorageMongo(m *mongo.Database, userID primitive.ObjectID, x int, y int, addVolume float64) bool {
	var storage StorageMongo

	res := m.Collection("storages").FindOne(context.TODO(), bson.M{"userId": userID, "x": x, "y": y})
	if res.Err() != nil {
		log.Println("Can't get storages: " + res.Err().Error())
		return false
	}
	err := res.Decode(&storage)
	if err != nil {
		log.Println("Can't decode storages: " + err.Error())
		return false
	}
	return storage.VolumeMax >= storage.VolumeOccupied+addVolume+GetDestinationVolumeMongo(m, userID, x, y)
}

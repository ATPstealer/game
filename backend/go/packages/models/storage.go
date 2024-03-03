package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

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

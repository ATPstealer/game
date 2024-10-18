package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Storage struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
	VolumeOccupied float64            `json:"volumeOccupied" bson:"volumeOccupied" validate:"required"`
	VolumeMax      float64            `json:"volumeMax" bson:"volumeMax" validate:"required"`
	X              int                `json:"x" bson:"x" validate:"required"`
	Y              int                `json:"y" bson:"y" validate:"required"`
} // @name storage

func GetAllStorages(m *mongo.Database) ([]Storage, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var allStorages []Storage
	cursor, err := m.Collection("storages").Find(ctx, bson.M{})
	if err != nil {
		return allStorages, err
	}

	err = cursor.All(ctx, &allStorages)
	return allStorages, err
}

func GetMyStorages(m *mongo.Database, userId primitive.ObjectID) ([]Storage, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var allStorages []Storage
	cursor, err := m.Collection("storages").Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return allStorages, err
	}

	err = cursor.All(ctx, &allStorages)
	return allStorages, err
}

func CheckEnoughStorage(m *mongo.Database, userId primitive.ObjectID, x int, y int, addVolume float64) bool {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var storage Storage
	res := m.Collection("storages").FindOne(ctx, bson.M{"userId": userId, "x": x, "y": y})
	if res.Err() != nil {
		log.Println("Can't get storages: " + res.Err().Error())
		return false
	}
	err := res.Decode(&storage)
	if err != nil {
		log.Println("Can't decode storages: " + err.Error())
		return false
	}
	return storage.VolumeMax >= storage.VolumeOccupied+addVolume+GetDestinationVolume(m, userId, x, y)
}

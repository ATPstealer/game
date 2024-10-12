package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ResourceType struct {
	Id         uint    `json:"id" bson:"id"`
	Name       string  `json:"name" bson:"name"`
	Volume     float64 `json:"volume" bson:"volume"` // m3
	Weight     float64 `json:"weight" bson:"weight"` // kg
	Demand     float64 `json:"demand" bson:"demand"`
	StoreGroup string  `json:"storeGroup" bson:"storeGroup"`
} // @name resourceType

func GetAllResourceTypes(m *mongo.Database) ([]ResourceType, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var resourceTypes []ResourceType
	cursor, err := m.Collection("resourceTypes").Find(ctx, bson.M{})
	if err != nil {
		return resourceTypes, err
	}

	err = cursor.All(ctx, &resourceTypes)
	return resourceTypes, err
}

func GetResourceTypesByID(m *mongo.Database, typeId uint) (ResourceType, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var resourceType ResourceType
	err := m.Collection("resourceTypes").FindOne(ctx, bson.M{"id": typeId}).Decode(&resourceType)
	if err != nil {
		log.Println("Can't get resource type: " + err.Error())
	}
	return resourceType, err
}

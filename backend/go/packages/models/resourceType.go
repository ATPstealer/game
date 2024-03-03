package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type ResourceTypeMongo struct {
	ID         uint    `json:"id" bson:"id"`
	Name       string  `json:"name" bson:"name"`
	Volume     float64 `json:"volume" bson:"volume"` // m3
	Weight     float64 `json:"weight" bson:"weight"` // kg
	Demand     float64 `json:"demand" bson:"demand"`
	StoreGroup string  `json:"storeGroup" bson:"storeGroup"`
}

func GetAllResourceTypesMongo(m *mongo.Database) ([]ResourceTypeMongo, error) {
	var resourceTypes []ResourceTypeMongo
	cursor, err := m.Collection("resourceTypes").Find(context.Background(), bson.M{})
	if err != nil {
		return resourceTypes, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.TODO(), &resourceTypes)
	return resourceTypes, err
}

func GetResourceTypesByIDMongo(m *mongo.Database, typeID uint) (ResourceTypeMongo, error) {
	var resourceType ResourceTypeMongo
	err := m.Collection("resourceTypes").FindOne(context.Background(), bson.M{"id": typeID}).Decode(&resourceType)
	if err != nil {
		log.Println("Can't get resource type: " + err.Error())
	}
	return resourceType, err
}

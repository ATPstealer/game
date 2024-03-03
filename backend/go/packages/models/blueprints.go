package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Blueprint struct {
	ID                uint             `json:"id" bson:"id"`
	Name              string           `json:"name" bson:"name"`
	ProducedResources []ResourceAmount `json:"producedResources" bson:"producedResources"`
	UsedResources     []ResourceAmount `json:"usedResources" bson:"usedResources"`
	ProducedInID      uint             `json:"producedInId" bson:"producedInId"`
	ProductionTime    time.Duration    `json:"productionTime" bson:"productionTime"`
}

type ResourceAmount struct {
	ResourceID uint    `json:"resourceId" bson:"resourceId"`
	Amount     float64 `json:"amount" bson:"amount"`
}

func GetBlueprints(m *mongo.Database, blueprintID uint) ([]Blueprint, error) {
	var blueprints []Blueprint
	filter := bson.M{}
	if blueprintID != 0 {
		filter["id"] = blueprintID
	}
	cursor, err := m.Collection("blueprints").Find(context.TODO(), filter,
		options.Find().SetSort(bson.D{{"blueprintId", 1}}))
	if err != nil {
		log.Println("Can't get blueprints: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &blueprints)
	return blueprints, err
}

func GetBlueprintByID(m *mongo.Database, blueprintID uint) (Blueprint, error) {
	var blueprint Blueprint
	err := m.Collection("blueprints").FindOne(context.TODO(),
		bson.M{"id": blueprintID}).Decode(&blueprint)
	if err != nil {
		log.Println("Can't get blueprint by ID: " + err.Error())
	}
	return blueprint, err
}

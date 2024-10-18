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
	Id                uint             `json:"id" bson:"id"`
	Name              string           `json:"name" bson:"name"`
	ProducedResources []ResourceAmount `json:"producedResources" bson:"producedResources"`
	UsedResources     []ResourceAmount `json:"usedResources" bson:"usedResources"`
	ProducedInId      uint             `json:"producedInId" bson:"producedInId"`
	ProductionTime    time.Duration    `json:"productionTime" bson:"productionTime"`
} // @name blueprint

type ResourceAmount struct {
	ResourceId uint    `json:"resourceId" bson:"resourceId"`
	Amount     float64 `json:"amount" bson:"amount"`
} // @name resourceAmount

func GetBlueprints(m *mongo.Database, blueprintId uint) ([]Blueprint, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var blueprints []Blueprint
	filter := bson.M{}
	if blueprintId != 0 {
		filter["id"] = blueprintId
	}
	cursor, err := m.Collection("blueprints").Find(ctx, filter,
		options.Find().SetSort(bson.D{{"id", 1}}))
	if err != nil {
		log.Println("Can't get blueprints: " + err.Error())
		return nil, err
	}

	err = cursor.All(ctx, &blueprints)
	return blueprints, err
}

func GetBlueprintById(m *mongo.Database, blueprintId uint) (Blueprint, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var blueprint Blueprint
	err := m.Collection("blueprints").FindOne(ctx, bson.M{"id": blueprintId}).Decode(&blueprint)
	if err != nil {
		log.Println("Can't get blueprint by Id: " + err.Error())
	}
	return blueprint, err
}

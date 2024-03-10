package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type EvolutionPrice struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	PriceAverage   float64            `json:"priceAverage" bson:"priceAverage"`
	Demand         float64            `json:"demand" bson:"demand"`
	SellSum        int                `json:"sellSum" bson:"sellSum"`
	RevenueSum     float64            `json:"revenueSum" bson:"revenueSum"`
	SpendMax       float64            `json:"spendMax" bson:"spendMax"`
}

func GetAllEvolutionPrices(m *mongo.Database) ([]EvolutionPrice, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var evolutionPrices []EvolutionPrice
	cursor, err := m.Collection("evolutionPrices").Find(ctx, bson.M{})
	if err != nil {
		return evolutionPrices, err
	}

	err = cursor.All(ctx, &evolutionPrices)
	return evolutionPrices, err
}

func GetEvolutionPrices(m *mongo.Database, x *int, y *int) ([]EvolutionPrice, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	if x != nil {
		filter = append(filter, bson.E{Key: "x", Value: *x})
	}
	if y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *y})
	}

	var evolutionPrices []EvolutionPrice
	cursor, err := m.Collection("evolutionPrices").Find(ctx, filter)
	if err != nil {
		return evolutionPrices, err
	}

	err = cursor.All(ctx, &evolutionPrices)
	return evolutionPrices, err
}

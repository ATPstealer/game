package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EvolutionPrice struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	PriceAverage   float64            `json:"priceAverage" bson:"priceAverage"`
	Demand         float64            `json:"demand" bson:"demand"`
	SellSum        int                `json:"sellSum" bson:"sellSum"`
	RevenueSum     float64            `json:"revenueSum" bson:"revenueSum"`
	SpendMax       float64            `json:"spendMax" bson:"spendMax"`
}

func GetAllEvolutionPrices(m *mongo.Database) ([]EvolutionPrice, error) {
	var evolutionPrices []EvolutionPrice
	cursor, err := m.Collection("evolutionPrices").Find(context.TODO(), bson.M{})
	if err != nil {
		return evolutionPrices, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.TODO(), &evolutionPrices)
	return evolutionPrices, err
}

func GetEvolutionPrices(m *mongo.Database, x *int, y *int) ([]EvolutionPrice, error) {
	filter := bson.D{}
	if x != nil {
		filter = append(filter, bson.E{Key: "x", Value: *x})
	}
	if y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *y})
	}

	var evolutionPrices []EvolutionPrice
	cursor, err := m.Collection("evolutionPrices").Find(context.TODO(), filter)
	if err != nil {
		return evolutionPrices, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.TODO(), &evolutionPrices)
	return evolutionPrices, err
}

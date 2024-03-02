package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
)

type EvolutionPrice struct {
	gorm.Model
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	PriceAverage   float64 `json:"priceAverage"`
	Demand         float64 `json:"demand"`
	SellSum        int     `json:"sellSum"`
	RevenueSum     float64 `json:"revenueSum"`
	SpendMax       float64 `json:"spendMax"`
}

func GetAllEvolutionPrices(db *gorm.DB) ([]EvolutionPrice, error) {
	var evolutionPrices []EvolutionPrice
	res := db.Model(&EvolutionPrice{}).Find(&evolutionPrices)
	if res.Error != nil {
		return nil, res.Error
	}
	return evolutionPrices, nil
}

type EvolutionPriceResult struct {
	ID             uint    `json:"id"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	PriceAverage   float64 `json:"priceAverage"`
}

func GetEvolutionPrices(db *gorm.DB, x *int, y *int) ([]EvolutionPriceResult, error) {
	var evolutionPrice []EvolutionPriceResult
	query := db.Model(&EvolutionPrice{})
	if x != nil {
		query = query.Where("x = ?", *x)
	}
	if y != nil {
		query = query.Where("y = ?", *y)
	}
	res := query.
		Select("id", "x", "y", "resource_type_id", "price_average").
		Scan(&evolutionPrice)

	if res.Error != nil {
		log.Println("Can't get evolution prices: " + res.Error.Error())
	}
	return evolutionPrice, res.Error
}

// mongo

type EvolutionPriceMongo struct {
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

func GetAllEvolutionPricesMongo(m *mongo.Database) ([]EvolutionPriceMongo, error) {
	var evolutionPrices []EvolutionPriceMongo
	cursor, err := m.Collection("evolutionPrices").Find(context.TODO(), bson.M{})
	if err != nil {
		return evolutionPrices, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.TODO(), &evolutionPrices)
	return evolutionPrices, err
}

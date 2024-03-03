package evolution

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
)

func ResetStats(db *gorm.DB) {
	resetStores(db)
	resetEvolutionPrices(db)
}

func resetStores(db *gorm.DB) {
	var storeGoods []models.StoreGoods
	db.Find(&storeGoods)
	for sgIndex := range storeGoods {
		storeGoods[sgIndex].Revenue = 0
		storeGoods[sgIndex].SellSum = 0
	}
	db.Save(storeGoods)
}

func resetEvolutionPrices(db *gorm.DB) {
	var evolutionPrices []models.EvolutionPrice
	db.Find(&evolutionPrices)
	for epIndex := range evolutionPrices {
		evolutionPrices[epIndex].SellSum = 0
		evolutionPrices[epIndex].RevenueSum = 0
	}
	db.Save(evolutionPrices)
}

//mongo

func ResetStatsMongo(m *mongo.Database) {
	resetStoresMongo(m)
	resetEvolutionPricesMongo(m)
}

func resetStoresMongo(m *mongo.Database) {
	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"sellSum": 0,
			"revenue": 0,
		},
	}
	_, err := m.Collection("storeGoods").UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
	}
}

func resetEvolutionPricesMongo(m *mongo.Database) {
	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"sellSum":    0,
			"revenueSum": 0,
		},
	}
	_, err := m.Collection("evolutionPriceS").UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
	}
}

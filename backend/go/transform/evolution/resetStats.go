package evolution

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func ResetStats(m *mongo.Database) {
	resetStores(m)
	resetEvolutionPrices(m)
}

func resetStores(m *mongo.Database) {
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

func resetEvolutionPrices(m *mongo.Database) {
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

package evolution

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

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

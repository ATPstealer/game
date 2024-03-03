package evolution

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func ResetStats(m *mongo.Database) {
	resetStores(m)
	resetEvolutionPrices(m)
}

func resetStores(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"sellSum": 0,
			"revenue": 0,
		},
	}
	_, err := m.Collection("storeGoods").UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
}

func resetEvolutionPrices(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"sellSum":    0,
			"revenueSum": 0,
		},
	}
	_, err := m.Collection("evolutionPriceS").UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
}

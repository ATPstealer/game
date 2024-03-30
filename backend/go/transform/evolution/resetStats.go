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
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"goods", bson.D{{"$ne", nil}}}}
	update := bson.D{
		{"$set",
			bson.D{
				{"goods.$[].sellSum", 0},
				{"goods.$[].revenue", 0},
			},
		},
	}
	if _, err := m.Collection("buildings").UpdateOne(ctx, filter, update); err != nil {
		log.Println("Can't reset store Goods stats: " + err.Error())
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
	_, err := m.Collection("evolutionPrices").UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
}

package db

import (
	"backend/packages/cfg"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	M *mongo.Database
)

func MongoConnect(config cfg.Vars) {
	timeoutDuration := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	clientOptions := options.Client().ApplyURI(
		"mongodb://" + config.MongoUser + ":" + config.MongoPassword + "@" + config.MongoHost + ":" + config.MongoPort)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	M = client.Database(config.MongoDatabase)
}

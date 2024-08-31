package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetAllReadyStorages(m *mongo.Database) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var readyStorages []Building

	filter := bson.M{"status": ReadyStatus, "onStrike": false, "typeId": 1} // TODO: надо конечно не так искать!
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return readyStorages, err
	}

	err = cursor.All(ctx, &readyStorages)
	return readyStorages, err
}

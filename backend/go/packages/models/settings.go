package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Settings struct {
	Key   string  `bson:"key" json:"key" validate:"required"`
	Value float64 `bson:"value" json:"value" validate:"required"`
} // @name settings

func GetSettings(m *mongo.Database) (map[string]float64, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	cursor, err := m.Collection("settings").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var settingsMongo []Settings
	if err = cursor.All(ctx, &settingsMongo); err != nil {
		return nil, err
	}

	settingsMap := make(map[string]float64)
	for _, setting := range settingsMongo {
		settingsMap[setting.Key] = setting.Value
	}

	return settingsMap, nil
}

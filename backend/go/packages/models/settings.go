package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

func GetSettingsMap(db *gorm.DB) map[string]float64 {
	var settings []Settings
	db.Find(&settings)
	settingsMap := make(map[string]float64)
	for _, set := range settings {
		settingsMap[set.Key] = set.Value
	}
	return settingsMap
}

// MONGO

type SettingsMongo struct {
	Key   string  `bson:"key" json:"key"`
	Value float64 `bson:"value" json:"value"`
}

func GetSettingsMongo(m *mongo.Database) (map[string]float64, error) {
	filter := bson.D{}
	cursor, err := m.Collection("settings").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var settingsMongo []SettingsMongo
	if err = cursor.All(context.TODO(), &settingsMongo); err != nil {
		return nil, err
	}

	settingsMap := make(map[string]float64)
	for _, setting := range settingsMongo {
		settingsMap[setting.Key] = setting.Value
	}

	return settingsMap, nil
}

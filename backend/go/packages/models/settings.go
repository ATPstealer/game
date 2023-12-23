package models

import (
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	Key   string  `json:"key"`
	Value float32 `json:"value"`
}

func GetSettingsMap(db *gorm.DB) map[string]float32 {
	var settings []Settings
	db.Find(&settings)
	settingsMap := make(map[string]float32)
	for _, set := range settings {
		settingsMap[set.Key] = set.Value
	}
	return settingsMap
}

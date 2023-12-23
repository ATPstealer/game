package models

import "gorm.io/gorm"

type EvolutionPrice struct {
	gorm.Model
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	PriceAverage   float32 `json:"priceAverage"`
	Demand         float32 `json:"demand"`
	SellSum        int     `json:"sellSum"`
	RevenueSum     float32 `json:"revenueSum"`
	SpendMax       float32 `json:"spendMax"`
}

func GetAllEvolutionPrices(db *gorm.DB) ([]EvolutionPrice, error) {
	var evolutionPrices []EvolutionPrice
	res := db.Model(&EvolutionPrice{}).Find(&evolutionPrices)
	if res.Error != nil {
		return nil, res.Error
	}
	return evolutionPrices, nil
}

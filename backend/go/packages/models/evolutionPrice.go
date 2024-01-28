package models

import (
	"gorm.io/gorm"
	"log"
)

type EvolutionPrice struct {
	gorm.Model
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	PriceAverage   float64 `json:"priceAverage"`
	Demand         float64 `json:"demand"`
	SellSum        int     `json:"sellSum"`
	RevenueSum     float64 `json:"revenueSum"`
	SpendMax       float64 `json:"spendMax"`
}

func GetAllEvolutionPrices(db *gorm.DB) ([]EvolutionPrice, error) {
	var evolutionPrices []EvolutionPrice
	res := db.Model(&EvolutionPrice{}).Find(&evolutionPrices)
	if res.Error != nil {
		return nil, res.Error
	}
	return evolutionPrices, nil
}

type EvolutionPriceResult struct {
	ID             uint    `json:"id"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	PriceAverage   float64 `json:"priceAverage"`
}

func GetEvolutionPrices(db *gorm.DB, x *int, y *int) ([]EvolutionPriceResult, error) {
	var evolutionPrice []EvolutionPriceResult
	query := db.Model(&EvolutionPrice{})
	if x != nil {
		query = query.Where("x = ?", *x)
	}
	if y != nil {
		query = query.Where("y = ?", *y)
	}
	res := query.
		Select("id", "x", "y", "resource_type_id", "price_average").
		Scan(&evolutionPrice)

	if res.Error != nil {
		log.Println("Can't get evolution prices: " + res.Error.Error())
	}
	return evolutionPrice, res.Error
}

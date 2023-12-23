package evolution

import (
	"backend/packages/models"
	"gorm.io/gorm"
)

func ResetStats(db *gorm.DB) {
	resetStores(db)
	resetEvolutionPrices(db)
}

func resetStores(db *gorm.DB) {
	var storeGoods []models.StoreGoods
	db.Find(&storeGoods)
	for sgIndex := range storeGoods {
		storeGoods[sgIndex].Revenue = 0
		storeGoods[sgIndex].SellSum = 0
	}
	db.Save(storeGoods)
}

func resetEvolutionPrices(db *gorm.DB) {
	var evolutionPrices []models.EvolutionPrice
	db.Find(&evolutionPrices)
	for epIndex := range evolutionPrices {
		evolutionPrices[epIndex].SellSum = 0
		evolutionPrices[epIndex].RevenueSum = 0
	}
	db.Save(evolutionPrices)
}

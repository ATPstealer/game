package evolution

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func CellSpendMaxMongo(m *mongo.Database) {
	cells, err := models.GetAllCellsMongo(m)
	if err != nil {
		log.Println("Can't get cells: " + err.Error())
		return
	}
	evolutionPrices, err := models.GetAllEvolutionPricesMongo(m)
	if err != nil {
		log.Println("Can't get evolution prices: " + err.Error())
		return
	}
	var cellPrices []models.EvolutionPriceMongo
	var newCellPrices []models.EvolutionPriceMongo
	for _, cell := range cells {
		cellPrices = findCellPricesMongo(&evolutionPrices, cell.X, cell.Y)
		sortByTurnoverMongo(&cellPrices)
		countSpendMaxMongo(&cellPrices, cell.CivilSavings*cell.SpendRate)
		newCellPrices = append(newCellPrices, cellPrices...)
	}
	saveSpend(m, &newCellPrices)
}

func findCellPricesMongo(evolutionPrices *[]models.EvolutionPriceMongo, x int, y int) []models.EvolutionPriceMongo {
	var cellPrices []models.EvolutionPriceMongo
	for _, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y {
			cellPrices = append(cellPrices, evolutionPrice)
		}
	}
	return cellPrices
}

func sortByTurnoverMongo(cellPrices *[]models.EvolutionPriceMongo) {
	for j := 0; j < len(*cellPrices); j++ {
		for k := j + 1; k < len(*cellPrices); k++ {
			if ((*cellPrices)[j].PriceAverage * (*cellPrices)[j].Demand) > ((*cellPrices)[k].PriceAverage * (*cellPrices)[k].Demand) {
				(*cellPrices)[j], (*cellPrices)[k] = (*cellPrices)[k], (*cellPrices)[j]
			}
		}
	}
}

func countSpendMaxMongo(cellPrices *[]models.EvolutionPriceMongo, moneyAvailable float64) {
	arrayLen := len(*cellPrices)
	var consumptionLevel float64
	for i := 0; i < arrayLen; i++ {
		var lastSpend float64
		if i == 0 {
			lastSpend = 0
		} else {
			lastSpend = (*cellPrices)[i-1].PriceAverage * (*cellPrices)[i-1].Demand
		}
		consumptionLevel = ((*cellPrices)[i].Demand*((*cellPrices)[i].PriceAverage) - lastSpend) * float64(arrayLen-i)
		if consumptionLevel >= moneyAvailable && i != 0 { // Подумать правильно ли это. ( i != 0 ) Раньше не было. Go вылетает из range на -1
			for j := i; j < arrayLen; j++ {
				(*cellPrices)[j].SpendMax = (*cellPrices)[i-1].Demand*(*cellPrices)[i-1].PriceAverage + moneyAvailable/float64(arrayLen-i)
			}
			return
		} else {
			(*cellPrices)[i].SpendMax = (*cellPrices)[i].Demand * (*cellPrices)[i].PriceAverage
			moneyAvailable = moneyAvailable - consumptionLevel
		}
	}
}

func saveSpend(m *mongo.Database, evolutionPrices *[]models.EvolutionPriceMongo) {
	for _, evolutionPrice := range *evolutionPrices {
		filter := bson.M{"_id": evolutionPrice.ID}
		update := bson.M{
			"$set": bson.M{
				"spendMax": evolutionPrice.SpendMax,
			},
		}

		_, err := m.Collection("evolutionPrices").UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}
	}
}

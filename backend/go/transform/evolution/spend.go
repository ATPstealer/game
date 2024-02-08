package evolution

import (
	"backend/packages/models"
	"gorm.io/gorm"
	"log"
)

func CellSpendMax(db *gorm.DB) {
	cells := models.GetAllCells(db)
	evolutionPrices, err := models.GetAllEvolutionPrices(db)
	if err != nil {
		log.Fatalln(err)
	}
	var cellPrices []models.EvolutionPrice
	var newCellPrices []models.EvolutionPrice
	for _, cell := range cells {
		cellPrices = findCellPrices(&evolutionPrices, cell.X, cell.Y)
		sortByTurnover(&cellPrices)
		countSpendMax(&cellPrices, cell.CivilSavings*cell.SpendRate)
		newCellPrices = append(newCellPrices, cellPrices...)
	}
	db.Save(&newCellPrices)
}

func findCellPrices(evolutionPrices *[]models.EvolutionPrice, x int, y int) []models.EvolutionPrice {
	var cellPrices []models.EvolutionPrice
	for _, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y {
			cellPrices = append(cellPrices, evolutionPrice)
		}
	}
	return cellPrices
}

func sortByTurnover(cellPrices *[]models.EvolutionPrice) {
	for j := 0; j < len(*cellPrices); j++ {
		for k := j + 1; k < len(*cellPrices); k++ {
			if ((*cellPrices)[j].PriceAverage * (*cellPrices)[j].Demand) > ((*cellPrices)[k].PriceAverage * (*cellPrices)[k].Demand) {
				(*cellPrices)[j], (*cellPrices)[k] = (*cellPrices)[k], (*cellPrices)[j]
			}
		}
	}
}

func countSpendMax(cellPrices *[]models.EvolutionPrice, moneyAvailable float64) {
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
		if consumptionLevel >= moneyAvailable {
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

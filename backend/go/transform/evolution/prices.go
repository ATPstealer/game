package evolution

import (
	"backend/packages/models"
	"gorm.io/gorm"
)

func CellAveragePrices(db *gorm.DB) error {
	cells := models.GetAllCells(db)
	allResourceTypes, err := models.GetAllResourceTypes(db)
	if err != nil {
		return err
	}
	storeGoods, err := models.GetAllStoreGoods(db)
	if err != nil {
		return err
	}
	evolutionPrices, err := models.GetAllEvolutionPrices(db)
	if err != nil {
		return err
	}
	for _, cell := range cells {
		for _, resourceType := range allResourceTypes {
			if resourceType.StoreGroup == "-" {
				continue
			}
			demand := resourceType.Demand * cell.Population
			cellGoods := findCellGoods(cell.X, cell.Y, resourceType.ID, &storeGoods)
			averagePrice := getAveragePrice(demand, cellGoods)
			addOrChangeEvolutionPrice(&evolutionPrices, cell.X, cell.Y, resourceType.ID, averagePrice, demand)
		}
	}
	db.Save(&evolutionPrices)
	return nil
}

func findCellGoods(x int, y int, resourceTypeID uint, storeGoods *[]models.StoreGoodsResult) []models.StoreGoodsResult {
	var cellGoods []models.StoreGoodsResult
	for _, sg := range *storeGoods {
		if sg.X == x && sg.Y == y && sg.ResourceTypeID == resourceTypeID {
			cellGoods = append(cellGoods, sg)
		}
	}
	return cellGoods
}

func getAveragePrice(demand float32, cellGoods []models.StoreGoodsResult) float32 {
	if len(cellGoods) == 0 || !goodsPriceExist(cellGoods) {
		return 0.01 // minimal Price for start selling
	}
	sort(&cellGoods)
	soldGoodsCount := float32(0)
	revenueCount := float32(0)
	for _, cg := range cellGoods {
		if soldGoodsCount+float32(cg.SellSum) <= demand {
			revenueCount += cg.Revenue
			soldGoodsCount += float32(cg.SellSum)
		} else {
			if soldGoodsCount+float32(cg.SellSum) == demand {
				break
			}
			revenueCount += (cg.Revenue / float32(cg.SellSum)) * (demand - soldGoodsCount)
			soldGoodsCount = demand
		}
	}
	return revenueCount / soldGoodsCount
}

func goodsPriceExist(cellGoods []models.StoreGoodsResult) bool {
	for _, cg := range cellGoods {
		if cg.SellSum != 0 {
			return true
		}
	}
	return false
}

func sort(cellGoods *[]models.StoreGoodsResult) {
	for j := 0; j < len(*cellGoods); j++ {
		for k := j + 1; k < len(*cellGoods); k++ {
			if (*cellGoods)[j].Revenue/float32((*cellGoods)[j].SellSum) > (*cellGoods)[k].Revenue/float32((*cellGoods)[k].SellSum) {
				(*cellGoods)[j], (*cellGoods)[k] = (*cellGoods)[k], (*cellGoods)[j]
			}
		}
	}
}

func addOrChangeEvolutionPrice(evolutionPrices *[]models.EvolutionPrice, x int, y int, resourceTypeID uint, averagePrice float32, demand float32) {
	for i, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y && evolutionPrice.ResourceTypeID == resourceTypeID {
			(*evolutionPrices)[i].PriceAverage = averagePrice
			(*evolutionPrices)[i].Demand = demand
			return
		}
	}
	*evolutionPrices = append(*evolutionPrices, models.EvolutionPrice{
		X:              x,
		Y:              y,
		ResourceTypeID: resourceTypeID,
		PriceAverage:   averagePrice,
		Demand:         demand,
		SellSum:        0,
		RevenueSum:     0,
	})
}

// TODO: reset count in stores

package evolution

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
)

func CellAveragePrices(db *gorm.DB) {
	cells := models.GetAllCells(db)
	allResourceTypes, err := models.GetAllResourceTypes(db)
	if err != nil {
		log.Fatalln(err)
	}
	storeGoods, err := models.GetAllStoreGoods(db)
	if err != nil {
		log.Fatalln(err)
	}
	evolutionPrices, err := models.GetAllEvolutionPrices(db)
	if err != nil {
		log.Fatalln(err)
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

func getAveragePrice(demand float64, cellGoods []models.StoreGoodsResult) float64 {
	if len(cellGoods) == 0 || !goodsPriceExist(cellGoods) {
		return 1 // Price for start selling
	}
	sort(&cellGoods)
	soldGoodsCount := float64(0)
	revenueCount := float64(0)
	for _, cg := range cellGoods {
		if soldGoodsCount+float64(cg.SellSum) <= demand {
			revenueCount += cg.Revenue
			soldGoodsCount += float64(cg.SellSum)
		} else {
			if soldGoodsCount+float64(cg.SellSum) == demand {
				break
			}
			revenueCount += (cg.Revenue / float64(cg.SellSum)) * (demand - soldGoodsCount)
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
			if (*cellGoods)[j].Revenue/float64((*cellGoods)[j].SellSum) > (*cellGoods)[k].Revenue/float64((*cellGoods)[k].SellSum) {
				(*cellGoods)[j], (*cellGoods)[k] = (*cellGoods)[k], (*cellGoods)[j]
			}
		}
	}
}

func addOrChangeEvolutionPrice(evolutionPrices *[]models.EvolutionPrice, x int, y int, resourceTypeID uint, averagePrice float64, demand float64) {
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

// mongo

func CellAveragePricesMongo(m *mongo.Database) {
	cells, err := models.GetAllCellsMongo(m)
	if err != nil {
		log.Fatalln(err)
	}

	allResourceTypes, err := models.GetAllResourceTypesMongo(m)
	if err != nil {
		log.Fatalln(err)
	}

	storeGoods, err := models.GetStoreGoodsWithDataMongo(m)
	if err != nil {
		log.Fatalln(err)
	}

	evolutionPrices, err := models.GetAllEvolutionPricesMongo(m)
	if err != nil {
		log.Fatalln(err)
	}

	for _, cell := range cells {
		for _, resourceType := range allResourceTypes {
			if resourceType.StoreGroup == "-" {
				continue
			}
			demand := resourceType.Demand * cell.Population
			cellGoods := findCellGoodsMongo(cell.X, cell.Y, resourceType.ID, &storeGoods)
			averagePrice := getAveragePriceMongo(demand, cellGoods)
			addOrChangeEvolutionPriceMongo(&evolutionPrices, cell.X, cell.Y, resourceType.ID, averagePrice, demand)
		}
	}
	saveEvolutionPrices(m, &evolutionPrices)
}

func findCellGoodsMongo(x int, y int, resourceTypeID uint, storeGoods *[]models.StoreGoodsWithDataMongo) []models.StoreGoodsWithDataMongo {
	var cellGoods []models.StoreGoodsWithDataMongo
	for _, sg := range *storeGoods {
		if sg.Building.X == x && sg.Building.Y == y && sg.ResourceTypeID == resourceTypeID {
			cellGoods = append(cellGoods, sg)
		}
	}
	return cellGoods
}

func getAveragePriceMongo(demand float64, cellGoods []models.StoreGoodsWithDataMongo) float64 {
	if len(cellGoods) == 0 || !goodsPriceExistMongo(cellGoods) {
		return 1 // Price for start selling
	}
	sortMongo(&cellGoods)
	soldGoodsCount := float64(0)
	revenueCount := float64(0)
	for _, cg := range cellGoods {
		if soldGoodsCount+float64(cg.SellSum) <= demand {
			revenueCount += cg.Revenue
			soldGoodsCount += float64(cg.SellSum)
		} else {
			if soldGoodsCount+float64(cg.SellSum) == demand {
				break
			}
			revenueCount += (cg.Revenue / float64(cg.SellSum)) * (demand - soldGoodsCount)
			soldGoodsCount = demand
		}
	}
	return revenueCount / soldGoodsCount
}

func goodsPriceExistMongo(cellGoods []models.StoreGoodsWithDataMongo) bool {
	for _, cg := range cellGoods {
		if cg.SellSum != 0 {
			return true
		}
	}
	return false
}

func sortMongo(cellGoods *[]models.StoreGoodsWithDataMongo) {
	for j := 0; j < len(*cellGoods); j++ {
		for k := j + 1; k < len(*cellGoods); k++ {
			if (*cellGoods)[j].Revenue/float64((*cellGoods)[j].SellSum) > (*cellGoods)[k].Revenue/float64((*cellGoods)[k].SellSum) {
				(*cellGoods)[j], (*cellGoods)[k] = (*cellGoods)[k], (*cellGoods)[j]
			}
		}
	}
}

func addOrChangeEvolutionPriceMongo(evolutionPrices *[]models.EvolutionPriceMongo, x int, y int, resourceTypeID uint, averagePrice float64, demand float64) {
	for i, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y && evolutionPrice.ResourceTypeID == resourceTypeID {
			(*evolutionPrices)[i].PriceAverage = averagePrice
			(*evolutionPrices)[i].Demand = demand
			return
		}
	}
	*evolutionPrices = append(*evolutionPrices, models.EvolutionPriceMongo{
		X:              x,
		Y:              y,
		ResourceTypeID: resourceTypeID,
		PriceAverage:   averagePrice,
		Demand:         demand,
		SellSum:        0,
		RevenueSum:     0,
	})
}

func saveEvolutionPrices(m *mongo.Database, evolutionPrices *[]models.EvolutionPriceMongo) {
	for _, price := range *evolutionPrices {
		filter := bson.M{"x": price.X, "y": price.Y, "resourceTypeID": price.ResourceTypeID}
		update := bson.M{
			"$set": bson.M{
				"priceAverage": price.PriceAverage,
				"demand":       price.Demand,
			},
			"$setOnInsert": bson.M{
				"x":              price.X,
				"y":              price.Y,
				"ResourceTypeID": price.ResourceTypeID,
			},
		}

		_, err := m.Collection("evolutionPrices").UpdateOne(context.TODO(), filter, update,
			options.Update().SetUpsert(true))
		if err != nil {
			log.Println(err)
		}
	}
}

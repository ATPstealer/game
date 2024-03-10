package evolution

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func CellAveragePrices(m *mongo.Database) {
	cells, err := models.GetAllCells(m)
	if err != nil {
		log.Fatalln(err)
	}

	allResourceTypes, err := models.GetAllResourceTypes(m)
	if err != nil {
		log.Fatalln(err)
	}

	storeGoods, err := models.GetAllStoreGoodsWithData(m)
	if err != nil {
		log.Fatalln(err)
	}

	evolutionPrices, err := models.GetAllEvolutionPrices(m)
	if err != nil {
		log.Fatalln(err)
	}

	for _, cell := range cells {
		for _, resourceType := range allResourceTypes {
			if resourceType.StoreGroup == "-" {
				continue
			}
			demand := resourceType.Demand * cell.Population
			cellGoods := findCellGoods(cell.X, cell.Y, resourceType.Id, &storeGoods)
			averagePrice := getAveragePrice(demand, cellGoods)
			addOrChangeEvolutionPrice(&evolutionPrices, cell.X, cell.Y, resourceType.Id, averagePrice, demand)
		}
	}
	SaveEvolutionPrices(m, &evolutionPrices)
}

func findCellGoods(x int, y int, resourceTypeID uint, storeGoods *[]models.StoreGoodsWithData) []models.StoreGoodsWithData {
	var cellGoods []models.StoreGoodsWithData
	for _, sg := range *storeGoods {
		if sg.Building.X == x && sg.Building.Y == y && sg.ResourceTypeId == resourceTypeID {
			cellGoods = append(cellGoods, sg)
		}
	}
	return cellGoods
}

func getAveragePrice(demand float64, cellGoods []models.StoreGoodsWithData) float64 {
	if len(cellGoods) == 0 || !goodsPriceExist(cellGoods) {
		return 1 // Price for start selling
	}
	sort(&cellGoods)
	soldGoodsCount := float64(0)
	revenueCount := float64(0)
	for _, cg := range cellGoods {
		if cg.Building.OnStrike {
			continue
		}
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

func goodsPriceExist(cellGoods []models.StoreGoodsWithData) bool {
	for _, cg := range cellGoods {
		if cg.SellSum != 0 {
			return true
		}
	}
	return false
}

func sort(cellGoods *[]models.StoreGoodsWithData) {
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
		if evolutionPrice.X == x && evolutionPrice.Y == y && evolutionPrice.ResourceTypeId == resourceTypeID {
			(*evolutionPrices)[i].PriceAverage = averagePrice
			(*evolutionPrices)[i].Demand = demand
			return
		}
	}

	*evolutionPrices = append(*evolutionPrices, models.EvolutionPrice{
		X:              x,
		Y:              y,
		ResourceTypeId: resourceTypeID,
		PriceAverage:   averagePrice,
		Demand:         demand,
		SellSum:        0,
		RevenueSum:     0,
	})
}

func SaveEvolutionPrices(m *mongo.Database, evolutionPrices *[]models.EvolutionPrice) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	for _, price := range *evolutionPrices {
		filter := bson.M{"x": price.X, "y": price.Y, "resourceTypeID": price.ResourceTypeId}
		update := bson.M{
			"$set": bson.M{
				"priceAverage": price.PriceAverage,
				"demand":       price.Demand,
				"sellSum":      price.SellSum,
				"revenueSum":   price.RevenueSum,
			},
			"$setOnInsert": bson.M{
				"x":              price.X,
				"y":              price.Y,
				"resourceTypeId": price.ResourceTypeId,
			},
		}

		_, err := m.Collection("evolutionPrices").UpdateOne(ctx, filter, update,
			options.Update().SetUpsert(true))
		if err != nil {
			log.Println(err)
		}
	}
}

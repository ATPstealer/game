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

	stores, err := models.GetBuildingsStores(m)
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
			cellGoods := findCellGoods(cell.X, cell.Y, resourceType.Id, &stores)
			if cell.X == 0 && cell.Y == 0 {
				log.Println(cellGoods, demand)

			}
			averagePrice := getAveragePrice(demand, cellGoods)
			addOrChangeEvolutionPrice(&evolutionPrices, cell.X, cell.Y, resourceType.Id, averagePrice, demand)
		}
	}
	SaveEvolutionPrices(m, &evolutionPrices)
}

func findCellGoods(x int, y int, resourceTypeID uint, stores *[]models.BuildingWithData) []models.Goods {
	var cellGoods []models.Goods
	for _, s := range *stores {
		if s.X != x || s.Y != y || s.Goods == nil {
			continue
		}
		for _, sg := range *s.Goods {
			if sg.ResourceTypeId == resourceTypeID {
				cellGoods = append(cellGoods, sg)
				break
			}
		}
	}
	return cellGoods
}

func getAveragePrice(demand float64, cellGoods []models.Goods) float64 {
	if len(cellGoods) == 0 || !goodsPriceExist(cellGoods) {
		return 1 // TODO: Можно в этом случае брать цены как средняя от соседей
	}
	sort(&cellGoods)
	soldGoodsCount := float64(0)
	revenueCount := float64(0)
	for _, cg := range cellGoods {
		if cg.Status == "OnStrike" {
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

func goodsPriceExist(cellGoods []models.Goods) bool {
	for _, cg := range cellGoods {
		if cg.SellSum != 0 {
			return true
		}
	}
	return false
}

func sort(cellGoods *[]models.Goods) {
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
		update := bson.M{"$set": price}
		_, err := m.Collection("evolutionPrices").UpdateOne(ctx, filter, update,
			options.Update().SetUpsert(true))
		if err != nil {
			log.Fatalln("Can't update Evolution Price cell", price.X, ":", price.Y, " resource ID ",
				price.ResourceTypeId, ": ", err)
		}
	}
}

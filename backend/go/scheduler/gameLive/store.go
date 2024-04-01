package gameLive

import (
	"backend/packages/models"
	"backend/transform/evolution"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"
)

// TODO: хотябы разбей это говно на функции
func StoreSell(m *mongo.Database) {
	buildingsGoods, err := models.GetBuildingsStores(m)
	if err != nil {
		log.Fatalln("Can't get Buildings Goods: " + err.Error())
	}

	evolutionPrices, err := models.GetAllEvolutionPrices(m)
	if err != nil {
		log.Fatalln("Can't get Evolution Prices: " + err.Error())
	}

	now := time.Now()

	for _, building := range buildingsGoods {
		for _, goods := range *building.Goods {
			if goods.Price == 0 {
				continue
			}
			if building.OnStrike {
				statusUpdate(m, building.Id, &goods, models.OnStrike)
				continue
			}
			if float64(goods.SellSum) >= building.BuildingType.Capacity*float64(building.Level*building.Square) {
				statusUpdate(m, building.Id, &goods, models.CapacityReached)
				continue
			}
			epIndex := findEvolutionPrice(&evolutionPrices, building.X, building.Y, goods.ResourceTypeId)
			if float64(evolutionPrices[epIndex].SellSum) >= evolutionPrices[epIndex].Demand {
				statusUpdate(m, building.Id, &goods, models.DemandSatisfied)
				continue
			}
			if evolutionPrices[epIndex].RevenueSum >= evolutionPrices[epIndex].SpendMax {
				statusUpdate(m, building.Id, &goods, models.SpendingLimitReached)
				continue
			}

			// Formula of selling pace
			workTime := now.Sub(goods.SellStarted).Seconds()
			storeCapacity := building.BuildingType.Capacity * float64(building.Workers) / float64(building.BuildingType.Workers) // square and level in Workers count
			daySells := daySellCalc(goods.Price, evolutionPrices[epIndex].PriceAverage, storeCapacity)

			if daySells < 1 {
				statusUpdate(m, building.Id, &goods, models.HighPrice)
				continue
			}

			oneSellTime := time.Second * time.Duration(24*60*60/daySells)
			sellCycles := int(daySells * workTime / (24 * 60 * 60))
			if sellCycles == 0 {
				continue
			}

			if !models.CheckEnoughResources(m, goods.ResourceTypeId, building.UserId, building.X, building.Y, float64(sellCycles)) {
				statusUpdate(m, building.Id, &goods, models.NotEnoughMinerals)
				if err := models.BuildingSetSellStarted(m, building.Id, goods.ResourceTypeId, now); err != nil {
					log.Println("Can't update Sell time: " + err.Error())
				}
				continue
			}

			if err := models.AddResource(m, goods.ResourceTypeId, building.UserId,
				building.X, building.Y, (-1)*float64(sellCycles)); err != nil {
				log.Println("Can't update resources: " + err.Error())
			}

			err = models.AddMoney(m, building.UserId, goods.Price*float64(sellCycles))
			if err != nil {
				log.Println(err.Error())
			}
			err = models.AddCivilSavings(m, building.X, building.Y, (-1)*float64(sellCycles)*goods.Price)
			if err != nil {
				log.Println(err.Error())
			}

			if err := models.BuildingGoodsStatsUpdate(m, building.Id, goods.ResourceTypeId,
				goods.SellSum+sellCycles, goods.Revenue+float64(sellCycles)*goods.Price); err != nil {
				log.Println("Can't update goods stats: " + err.Error())
			}

			if err := models.BuildingSetSellStarted(m, building.Id, goods.ResourceTypeId,
				goods.SellStarted.Add(time.Duration(sellCycles)*oneSellTime)); err != nil {
				log.Println("Can't update Sell time: " + err.Error())
			}

			statusUpdate(m, building.Id, &goods, models.Selling)

			evolutionPrices[epIndex].SellSum += sellCycles
			evolutionPrices[epIndex].RevenueSum += float64(sellCycles) * goods.Price
		}
	}

	evolution.SaveEvolutionPrices(m, &evolutionPrices)
}

func daySellCalc(price float64, priceAverage float64, capacity float64) float64 {
	if priceAverage == 0 {
		return 0
	}
	return capacity * 0.75 * priceAverage / price
}

func findEvolutionPrice(evolutionPrices *[]models.EvolutionPrice, x int, y int, resourceTypeID uint) int {
	for index, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y && evolutionPrice.ResourceTypeId == resourceTypeID {
			return index
		}
	}
	log.Fatalln(fmt.Sprintf("can't find evolutionPrice in %s:%s resource %s",
		strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(int(resourceTypeID))))
	return -1
}

func statusUpdate(m *mongo.Database, buildingId primitive.ObjectID, goods *models.Goods, status models.StoreGoodsStatus) {
	if goods.Status != status {
		if err := models.BuildingGoodsStatusUpdate(m, buildingId, goods.ResourceTypeId, status); err != nil {
			log.Println("Building ID ", buildingId, "Can't update Goods ", (*goods).ResourceTypeId, " status: ", err.Error())
		}
	}
}

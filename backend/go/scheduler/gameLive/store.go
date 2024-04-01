package gameLive

import (
	"backend/packages/models"
	"backend/transform/evolution"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"
)

// TODO: хотябы разбей это говно на функции
func StoreSell(m *mongo.Database) {
	buildingsGoods, err := models.GetBuildingsStores(m)
	log.Println(buildingsGoods)
	if err != nil {
		log.Println("Can't get Buildings Goods: " + err.Error())
		return
	}

	evolutionPrices, err := models.GetAllEvolutionPrices(m)
	if err != nil {
		log.Println("Can't get Evolution Prices: " + err.Error())
		return
	}

	now := time.Now()

	for _, building := range buildingsGoods {
		for _, goods := range *building.Goods {
			if goods.Price == 0 {
				continue
			}
			if building.OnStrike {
				if goods.Status != models.OnStrike {
					if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.OnStrike); err != nil {
						log.Println("Can't update Goods status: " + err.Error())
					}
				}
				continue
			}
			log.Println(goods)

			if float64(goods.SellSum) >= building.BuildingType.Capacity*float64(building.Level*building.Square) {
				if goods.Status != models.CapacityReached {
					if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.CapacityReached); err != nil {
						log.Println("Can't update Goods status: " + err.Error())
					}
				}
				continue
			}

			epIndex, err := findEvolutionPrice(&evolutionPrices, building.X, building.Y, goods.ResourceTypeId)
			if err != nil {
				log.Println(err)
			}

			if float64(evolutionPrices[epIndex].SellSum) >= evolutionPrices[epIndex].Demand {
				if goods.Status != models.DemandSatisfied {
					if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.DemandSatisfied); err != nil {
						log.Println("Can't update Goods status: " + err.Error())
					}
				}
				continue
			}

			if evolutionPrices[epIndex].RevenueSum >= evolutionPrices[epIndex].SpendMax {
				if goods.Status != models.SpendingLimitReached {
					if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.SpendingLimitReached); err != nil {
						log.Println("Can't update Goods status: " + err.Error())
					}
				}
				continue
			}

			// Formula of selling pace
			workTime := now.Sub(goods.SellStarted).Seconds()
			storeCapacity := building.BuildingType.Capacity * float64(building.Workers) / float64(building.BuildingType.Workers) // square and level in Workers count
			daySells := daySellCalc(goods.Price, evolutionPrices[epIndex].PriceAverage, storeCapacity)
			oneSellTime := time.Second * time.Duration(24*60*60/daySells) // TODO: проверить что здесь нормальные секунды

			if daySells == 0 {
				if goods.Status != models.HighPrice {
					if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.HighPrice); err != nil {
						log.Println("Can't update Goods status: " + err.Error())
					}
				}
				continue
			} else {
				if goods.Status != models.Selling {
					if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.Selling); err != nil {
						log.Println("Can't update Goods status: " + err.Error())
					}
				}
			}
			sellCycles := int(daySells * workTime / (24 * 60 * 60))
			if sellCycles == 0 {
				continue
			}

			if !models.CheckEnoughResources(m, goods.ResourceTypeId, building.UserId,
				building.X, building.Y, float64(sellCycles)) {
				if goods.Status != models.NotEnoughMinerals {
					if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.NotEnoughMinerals); err != nil {
						log.Println("Can't update Goods status: " + err.Error())
					}
				}
				if err := models.BuildingSetSellStarted(m, building.Id, goods.ResourceTypeId, now); err != nil {
					log.Println("Can't update Sell time: " + err.Error())
				}
				continue
			}

			log.Println(goods.ResourceTypeId, building.UserId,
				building.X, building.Y, (-1)*float64(sellCycles))
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

			log.Println("test")

			if err := models.BuildingGoodsStatsUpdate(m, building.Id, goods.ResourceTypeId,
				goods.SellSum+sellCycles, goods.Revenue+float64(sellCycles)*goods.Price); err != nil {
				log.Println("Can't update goods stats: " + err.Error())
			}

			if err := models.BuildingSetSellStarted(m, building.Id, goods.ResourceTypeId,
				goods.SellStarted.Add(time.Duration(sellCycles)*oneSellTime)); err != nil {
				log.Println("Can't update Sell time: " + err.Error())
			}

			if goods.Status != models.Selling {
				if err := models.BuildingGoodsStatusUpdate(m, building.Id, goods.ResourceTypeId, models.Selling); err != nil {
					log.Println("Can't update Goods status: " + err.Error())
				}
			}

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

func findEvolutionPrice(evolutionPrices *[]models.EvolutionPrice, x int, y int, resourceTypeID uint) (int, error) {
	for index, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y && evolutionPrice.ResourceTypeId == resourceTypeID {
			return index, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("can't find evolutionPrice in %s:%s resource %s", strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(int(resourceTypeID))))
}

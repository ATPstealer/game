package gameLive

import (
	"backend/packages/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

func StoreSell(db *gorm.DB) error {
	storeGoods, err := models.GetAllStoreGoods(db)
	if err != nil {
		return err
	}
	evolutionPrices, err := models.GetAllEvolutionPrices(db)
	if err != nil {
		return err
	}
	now := time.Now()

	for gIndex, goods := range storeGoods {
		if goods.Price == 0 {
			continue
		}

		epIndex, err := findEvolutionPrice(&evolutionPrices, goods.X, goods.Y, goods.ResourceTypeID)
		if err != nil {
			log.Println(err)
		}

		if float64(evolutionPrices[epIndex].SellSum) >= evolutionPrices[epIndex].Demand {
			storeGoods[gIndex].Status = models.DemandSatisfied
			continue
		}

		if evolutionPrices[epIndex].RevenueSum >= evolutionPrices[epIndex].SpendMax {
			storeGoods[gIndex].Status = models.SpendingLimitReached
			continue
		}

		// Formula of selling pace
		workTime := float64(now.Sub(*goods.SellStarted).Seconds())
		storeCapacity := float64(goods.Capacity * goods.Level * goods.Square)
		daySells := daySellCalc(goods.Price, evolutionPrices[epIndex].PriceAverage, storeCapacity)
		oneSellTime := time.Second * time.Duration(24*60*60/daySells)

		if daySells == 0 {
			storeGoods[gIndex].Status = models.HighPrice
			continue
		}
		sellCycles := int(daySells * workTime / (24 * 60 * 60))

		if !models.CheckEnoughResources(db, goods.ResourceTypeID, goods.UserID, goods.X, goods.Y, float64(sellCycles)) {
			storeGoods[gIndex].Status = models.NotEnoughMinerals
			storeGoods[gIndex].SellStarted = &now
			continue
		}

		err = models.AddResource(db, goods.ResourceTypeID, goods.UserID, goods.X, goods.Y, (-1)*float64(sellCycles))
		if err != nil {
			log.Println(err.Error())
		}

		err = models.AddMoney(db, goods.UserID, goods.Price*float64(sellCycles))
		if err != nil {
			log.Println(err.Error())
		}
		err = models.AddCivilSavings(db, goods.X, goods.Y, (-1)*float64(sellCycles)*goods.Price)
		if err != nil {
			log.Println(err.Error())
		}

		storeGoods[gIndex].SellSum += sellCycles
		storeGoods[gIndex].Revenue += float64(sellCycles) * goods.Price
		evolutionPrices[epIndex].SellSum += sellCycles
		evolutionPrices[epIndex].RevenueSum += float64(sellCycles) * goods.Price

		newWorkStarted := storeGoods[gIndex].SellStarted.Add(time.Duration(sellCycles) * oneSellTime)
		storeGoods[gIndex].SellStarted = &newWorkStarted
		storeGoods[gIndex].Status = models.Selling
	}

	saveStoreGoods(db, &storeGoods)
	db.Save(&evolutionPrices)
	return nil
}

func daySellCalc(price float64, priceAverage float64, capacity float64) float64 {
	if priceAverage == 0 {
		return 0
	}
	return capacity * 0.75 * priceAverage / price
}

func saveStoreGoods(db *gorm.DB, storeGoodsResult *[]models.StoreGoodsResult) {
	var storeGoods []models.StoreGoods
	for _, sgr := range *storeGoodsResult {
		storeGoods = append(storeGoods, models.StoreGoods{
			Model:          gorm.Model{ID: sgr.ID},
			BuildingID:     sgr.BuildingID,
			ResourceTypeID: sgr.ResourceTypeID,
			Price:          sgr.Price,
			SellSum:        sgr.SellSum,
			Revenue:        sgr.Revenue,
			SellStarted:    sgr.SellStarted,
			Status:         sgr.Status,
		})
	}
	db.Save(storeGoods)
}

func findEvolutionPrice(evolutionPrices *[]models.EvolutionPrice, x int, y int, resourceTypeID uint) (int, error) {
	for index, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y && evolutionPrice.ResourceTypeID == resourceTypeID {
			return index, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("can't find evolutionPrice in %s:%s resource %s", strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(int(resourceTypeID))))
}

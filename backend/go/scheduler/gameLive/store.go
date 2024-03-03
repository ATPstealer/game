package gameLive

import (
	"backend/packages/models"
	"backend/transform/evolution"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"
)

func StoreSellMongo(m *mongo.Database) {
	storeGoods, err := models.GetAllStoreGoodsWithDataMongo(m)
	if err != nil {
		log.Println("Can't get Store Goods: " + err.Error())
		return
	}

	evolutionPrices, err := models.GetAllEvolutionPricesMongo(m)
	if err != nil {
		log.Println("Can't get Evolution Prices: " + err.Error())
		return
	}

	now := time.Now()

	for gIndex, goods := range storeGoods {
		if goods.Price == 0 {
			continue
		}

		if float64(goods.SellSum) >= goods.BuildingType.Capacity*float64(goods.Building.Level*goods.Building.Square) {
			storeGoods[gIndex].Status = models.CapacityReached
			continue
		}

		epIndex, err := findEvolutionPriceMongo(&evolutionPrices, goods.Building.X, goods.Building.Y, goods.ResourceTypeID)
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
		workTime := now.Sub(goods.SellStarted).Seconds()
		storeCapacity := goods.BuildingType.Capacity * float64(goods.Building.Workers) / float64(goods.BuildingType.Workers) // square and level in Workers count
		daySells := daySellCalcMongo(goods.Price, evolutionPrices[epIndex].PriceAverage, storeCapacity)
		oneSellTime := time.Second * time.Duration(24*60*60/daySells)

		if daySells == 0 {
			storeGoods[gIndex].Status = models.HighPrice
			continue
		} else {
			storeGoods[gIndex].Status = models.Selling
		}
		sellCycles := int(daySells * workTime / (24 * 60 * 60))
		if sellCycles == 0 {
			continue
		}

		if !models.CheckEnoughResourcesMongo(m, goods.ResourceTypeID, goods.Building.UserID,
			goods.Building.X, goods.Building.Y, float64(sellCycles)) {
			storeGoods[gIndex].Status = models.NotEnoughMinerals
			storeGoods[gIndex].SellStarted = now
			continue
		}

		err = models.AddResourceMongo(m, goods.ResourceTypeID, goods.Building.UserID,
			goods.Building.X, goods.Building.Y, (-1)*float64(sellCycles))
		if err != nil {
			log.Println(err.Error())
		}

		err = models.AddMoneyMongo(m, goods.Building.UserID, goods.Price*float64(sellCycles))
		if err != nil {
			log.Println(err.Error())
		}
		err = models.AddCivilSavingsMongo(m, goods.Building.X, goods.Building.Y, (-1)*float64(sellCycles)*goods.Price)
		if err != nil {
			log.Println(err.Error())
		}

		storeGoods[gIndex].SellSum += sellCycles
		storeGoods[gIndex].Revenue += float64(sellCycles) * goods.Price
		evolutionPrices[epIndex].SellSum += sellCycles
		evolutionPrices[epIndex].RevenueSum += float64(sellCycles) * goods.Price

		storeGoods[gIndex].SellStarted = storeGoods[gIndex].SellStarted.Add(time.Duration(sellCycles) * oneSellTime)
		storeGoods[gIndex].Status = models.Selling
	}

	saveStoreGoodsMongo(m, &storeGoods)
	evolution.SaveEvolutionPrices(m, &evolutionPrices)
}

func daySellCalcMongo(price float64, priceAverage float64, capacity float64) float64 {
	if priceAverage == 0 {
		return 0
	}
	return capacity * 0.75 * priceAverage / price
}

func findEvolutionPriceMongo(evolutionPrices *[]models.EvolutionPriceMongo, x int, y int, resourceTypeID uint) (int, error) {
	for index, evolutionPrice := range *evolutionPrices {
		if evolutionPrice.X == x && evolutionPrice.Y == y && evolutionPrice.ResourceTypeID == resourceTypeID {
			return index, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("can't find evolutionPrice in %s:%s resource %s", strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(int(resourceTypeID))))
}

func saveStoreGoodsMongo(m *mongo.Database, storeGoods *[]models.StoreGoodsWithDataMongo) {
	for _, sg := range *storeGoods {
		filter := bson.M{"buildingId": sg.BuildingID, "resourceTypeId": sg.ResourceTypeID}
		update := bson.M{
			"$set": bson.M{
				"sellSum":     sg.SellSum,
				"revenue":     sg.Revenue,
				"sellStarted": sg.SellStarted,
				"status":      sg.Status,
			},
		}
		_, err := m.Collection("storeGoods").UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}
	}
}

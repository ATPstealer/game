package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type StoreGoodsStatus string

const (
	Selling              StoreGoodsStatus = "Selling"
	DemandSatisfied      StoreGoodsStatus = "DemandSatisfied"
	HighPrice            StoreGoodsStatus = "HighPrice"
	NotEnoughMinerals    StoreGoodsStatus = "NotEnoughMinerals"
	SpendingLimitReached StoreGoodsStatus = "SpendingLimitReached"
	CapacityReached      StoreGoodsStatus = "CapacityReached"
	OnStrike             StoreGoodsStatus = "OnStrike"
)

type StoreGoodsPayload struct {
	BuildingId     primitive.ObjectID `json:"buildingId"`
	ResourceTypeId uint               `json:"resourceTypeId"`
	Price          float64            `json:"price"`
}

func SetStoreGoods(m *mongo.Database, userId primitive.ObjectID, payload StoreGoodsPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		return err
	}
	if building.UserId != userId {
		return errors.New("this building don't belong you")
	}

	buildingType, err := GetBuildingTypeById(m, building.TypeId)
	if err != nil {
		return err
	}
	if buildingType.BuildingGroup != "Store" {
		return errors.New("this is not a store")
	}

	resourceType, err := GetResourceTypesByID(m, payload.ResourceTypeId)
	if err != nil {
		return err
	}
	if resourceType.StoreGroup != buildingType.BuildingSubGroup {
		return errors.New("can't sell here")
	}

	index := getIdPosition(building.Goods, payload.ResourceTypeId)
	if index == -1 {
		newGoodsPrice := Goods{
			ResourceTypeId: payload.ResourceTypeId,
			Price:          payload.Price,
			SellSum:        0,
			Revenue:        0,
			SellStarted:    time.Now(),
			Status:         Selling,
		}
		newGoodsArr := []Goods{newGoodsPrice}
		if building.Goods != nil {
			newGoodsArr = append(*building.Goods, newGoodsPrice)
		}
		building.Goods = &newGoodsArr
	} else {
		(*building.Goods)[index].Price = payload.Price
		(*building.Goods)[index].SellStarted = time.Now()
		(*building.Goods)[index].Status = Selling
	}

	_, err = m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": building.Id},
		bson.M{"$set": building})

	return err
}

func getIdPosition(goodsArr *[]Goods, typeId uint) int {
	if goodsArr == nil {
		return -1
	}
	for i, v := range *goodsArr {
		if v.ResourceTypeId == typeId {
			return i
		}
	}
	return -1
}

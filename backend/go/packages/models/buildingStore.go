package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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

type Goods struct {
	ResourceTypeId uint             `json:"resourceTypeId" bson:"resourceTypeId"`
	Price          float64          `json:"price" bson:"price"`
	SellSum        int              `json:"sellSum" bson:"sellSum"`
	Revenue        float64          `json:"revenue" bson:"revenue"`
	SellStarted    time.Time        `json:"sellStarted" bson:"sellStarted"`
	Status         StoreGoodsStatus `json:"status" bson:"status"`
}

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
		return errors.New("this building doesn't belong to you")
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

func GetBuildingsStores(m *mongo.Database) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{
		{"workEnd", bson.D{{"$gt", time.Now()}}},
		{"goods", bson.D{{"$ne", nil}}},
	}
	matchStage := bson.D{{"$match", filter}}

	lookupBuildingType := bson.D{{"$lookup", bson.D{
		{"from", "buildingTypes"},
		{"localField", "typeId"},
		{"foreignField", "id"},
		{"as", "buildingType"},
	}}}

	unwindBuildingType := bson.D{{"$unwind", bson.D{
		{"path", "$buildingType"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	pipeline := mongo.Pipeline{matchStage, lookupBuildingType, unwindBuildingType}
	cursor, err := m.Collection("buildings").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get productions: " + err.Error())
		return nil, err
	}

	var buildingWithData []BuildingWithData
	if err = cursor.All(ctx, &buildingWithData); err != nil {
		log.Println(err)
	}
	return buildingWithData, nil
}

func BuildingGoodsStatusUpdate(m *mongo.Database, buildingId primitive.ObjectID, resourceTypeId uint, status StoreGoodsStatus) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"_id", buildingId}}
	update := bson.D{{"$set", bson.D{{"goods.$[elem].status", status}}}}
	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"elem.resourceTypeId", resourceTypeId}},
		},
	})

	_, err := m.Collection("buildings").UpdateOne(ctx, filter, update, updateOpts)
	return err
}

func BuildingSetSellStarted(m *mongo.Database, buildingId primitive.ObjectID, resourceTypeId uint, timeStart time.Time) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"_id", buildingId}}
	update := bson.D{{"$set", bson.D{{"goods.$[elem].sellStarted", timeStart}}}}
	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"elem.resourceTypeId", resourceTypeId}},
		},
	})

	_, err := m.Collection("buildings").UpdateOne(ctx, filter, update, updateOpts)
	return err
}

func BuildingGoodsStatsUpdate(m *mongo.Database, buildingId primitive.ObjectID, resourceTypeId uint, sellSum int, revenue float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"_id", buildingId}}
	update := bson.D{{"$set",
		bson.D{
			{"goods.$[elem].sellSum", sellSum},
			{"goods.$[elem].revenue", revenue},
		},
	}}
	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"elem.resourceTypeId", resourceTypeId}},
		},
	})

	_, err := m.Collection("buildings").UpdateOne(ctx, filter, update, updateOpts)
	return err
}

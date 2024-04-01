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

type StoreGoods struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingId     primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	Price          float64            `json:"price" bson:"price"`
	SellSum        int                `json:"sellSum" bson:"sellSum"`
	Revenue        float64            `json:"revenue" bson:"revenue"`
	SellStarted    time.Time          `json:"sellStarted" bson:"sellStarted"`
	Status         StoreGoodsStatus   `json:"status" bson:"status"`
}

func GetStoreGoods(m *mongo.Database, buildingId primitive.ObjectID) ([]StoreGoods, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var storeGoods []StoreGoods
	filter := bson.M{}
	if buildingId != primitive.NilObjectID {
		filter["buildingId"] = buildingId
	}
	cursor, err := m.Collection("storeGoods").Find(ctx, filter)
	if err != nil {
		log.Println("Can't get store goods: " + err.Error())
		return nil, err
	}

	err = cursor.All(ctx, &storeGoods)
	return storeGoods, err
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

	_, err = m.Collection("storeGoods").UpdateOne(ctx,
		bson.M{
			"buildingId":     payload.BuildingId,
			"resourceTypeId": payload.ResourceTypeId,
		},
		bson.M{
			"$set": bson.M{
				"price":       payload.Price,
				"sellStarted": time.Now(),
			},
			"$setOnInsert": bson.M{
				"buildingID":     payload.BuildingId,
				"resourceTypeId": payload.ResourceTypeId,
				"sellSum":        0,
				"revenue":        0,
				"status":         Selling,
			},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	log.Println("asdasd")
	index := getIdPosition(building.Goods, payload.ResourceTypeId)
	log.Println(index)
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
	log.Println("6")

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

type StoreGoodsWithData struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingId     primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	Price          float64            `json:"price" bson:"price"`
	SellSum        int                `json:"sellSum" bson:"sellSum"`
	Revenue        float64            `json:"revenue" bson:"revenue"`
	SellStarted    time.Time          `json:"sellStarted" bson:"sellStarted"`
	Status         StoreGoodsStatus   `json:"status" bson:"status"`
	Building       Building           `json:"building" bson:"building"`
	BuildingType   BuildingType       `json:"buildingType" bson:"buildingType"`
}

func GetAllStoreGoodsWithData(m *mongo.Database) ([]StoreGoodsWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{}}
	matchStage := bson.D{{"$match", filter}}

	lookupBuilding := bson.D{{"$lookup", bson.D{
		{"from", "buildings"},
		{"localField", "buildingId"},
		{"foreignField", "_id"},
		{"as", "building"},
	}}}

	unwindBuilding := bson.D{{"$unwind", bson.D{
		{"path", "$building"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	lookupBuildingType := bson.D{{"$lookup", bson.D{
		{"from", "buildingTypes"},
		{"localField", "building.typeId"},
		{"foreignField", "id"},
		{"as", "buildingType"},
	}}}

	unwindBuildingType := bson.D{{"$unwind", bson.D{
		{"path", "$buildingType"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	pipeline := mongo.Pipeline{matchStage, lookupBuilding, unwindBuilding, lookupBuildingType, unwindBuildingType}
	cursor, err := m.Collection("storeGoods").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get store goods: " + err.Error())
		return nil, err
	}

	var storeGoods []StoreGoodsWithData
	if err = cursor.All(ctx, &storeGoods); err != nil {
		log.Println(err)
	}
	return storeGoods, nil
}

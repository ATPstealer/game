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
)

type StoreGoodsMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingID     primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	Price          float64            `json:"price" bson:"price"`
	SellSum        int                `json:"sellSum" bson:"sellSum"`
	Revenue        float64            `json:"revenue" bson:"revenue"`
	SellStarted    time.Time          `json:"sellStarted" bson:"sellStarted"`
	Status         StoreGoodsStatus   `json:"status" bson:"status"`
}

func GetStoreGoodsMongo(m *mongo.Database, buildingID primitive.ObjectID) ([]StoreGoodsMongo, error) {
	var storeGoods []StoreGoodsMongo
	filter := bson.M{}
	if buildingID != primitive.NilObjectID {
		filter["buildingId"] = buildingID
	}
	cursor, err := m.Collection("storeGoods").Find(context.TODO(), filter)
	if err != nil {
		log.Println("Can't get store goods: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &storeGoods)
	return storeGoods, err
}

type StoreGoodsPayloadMongo struct {
	BuildingID     primitive.ObjectID `json:"buildingId"`
	ResourceTypeID uint               `json:"resourceTypeId"`
	Price          float64            `json:"price"`
}

func SetStoreGoodsMongo(m *mongo.Database, userID primitive.ObjectID, payload StoreGoodsPayloadMongo) error {
	building, err := GetBuildingByIDMongo(m, payload.BuildingID)
	if err != nil {
		return err
	}
	if building.UserID != userID {
		return errors.New("this building don't belong you")
	}

	buildingType, err := GetBuildingTypeByIDMongo(m, building.TypeID)
	if err != nil {
		return err
	}
	if buildingType.BuildingGroup != "Store" {
		return errors.New("this is not a store")
	}

	resourceType, err := GetResourceTypesByIDMongo(m, payload.ResourceTypeID)
	if err != nil {
		return err
	}
	if resourceType.StoreGroup != buildingType.BuildingSubGroup {
		return errors.New("can't sell here")
	}

	_, err = m.Collection("storeGoods").UpdateOne(context.TODO(),
		bson.M{
			"buildingId":     payload.BuildingID,
			"resourceTypeId": payload.ResourceTypeID,
		},
		bson.M{
			"$set": bson.M{
				"price":       payload.Price,
				"sellStarted": time.Now(),
			},
			"$setOnInsert": bson.M{
				"buildingID":     payload.BuildingID,
				"resourceTypeId": payload.ResourceTypeID,
				"sellSum":        0,
				"revenue":        0,
				"status":         Selling,
			},
		},
		options.Update().SetUpsert(true))
	return err
}

type StoreGoodsWithDataMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingID     primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	Price          float64            `json:"price" bson:"price"`
	SellSum        int                `json:"sellSum" bson:"sellSum"`
	Revenue        float64            `json:"revenue" bson:"revenue"`
	SellStarted    time.Time          `json:"sellStarted" bson:"sellStarted"`
	Status         StoreGoodsStatus   `json:"status" bson:"status"`
	Building       BuildingMongo      `json:"building" bson:"building"`
	BuildingType   BuildingTypeMongo  `json:"buildingType" bson:"buildingType"`
}

func GetAllStoreGoodsWithDataMongo(m *mongo.Database) ([]StoreGoodsWithDataMongo, error) {
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
	cursor, err := m.Collection("storeGoods").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get store goods: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var storeGoods []StoreGoodsWithDataMongo
	if err = cursor.All(context.TODO(), &storeGoods); err != nil {
		log.Println(err)
	}
	return storeGoods, nil
}

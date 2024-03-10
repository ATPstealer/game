package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type BuildingStatus string

const (
	ConstructionStatus    BuildingStatus = "Construction"
	ReadyStatus           BuildingStatus = "Ready"
	ProductionStatus      BuildingStatus = "Production"
	ResourcesNeededStatus BuildingStatus = "ResourcesNeeded"
	StorageNeededStatus   BuildingStatus = "StorageNeeded"
)

type Building struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TypeId      uint               `json:"typeId" bson:"typeId"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`
	X           int                `json:"x" bson:"x"`
	Y           int                `json:"y" bson:"y"`
	Square      int                `json:"square" bson:"square"`
	Level       int                `json:"level" bson:"level"`
	Status      BuildingStatus     `json:"status" bson:"status"`
	WorkStarted time.Time          `json:"workStarted" bson:"workStarted"`
	WorkEnd     time.Time          `json:"workEnd" bson:"workEnd"`
	HiringNeeds int                `json:"hiringNeeds" bson:"hiringNeeds"`
	Salary      float64            `json:"salary" bson:"salary"`
	Workers     int                `json:"workers" bson:"workers"`
	OnStrike    bool               `json:"onStrike" bson:"onStrike"`
}

type ConstructBuildingPayload struct {
	TypeId uint `json:"typeId"`
	X      int  `json:"x"`
	Y      int  `json:"y"`
	Square int  `json:"square"`
}

func ConstructBuilding(m *mongo.Database, userId primitive.ObjectID, payload ConstructBuildingPayload) error {
	enoughLand, err := CheckEnoughLandForBuilding(m, userId, payload.Square, payload.X, payload.Y)
	if !enoughLand {
		return errors.New("not enough land")
	}
	if err != nil {
		return err
	}

	buildingType, err := GetBuildingTypeByID(m, payload.TypeId)
	if err != nil {
		return err
	}
	if !CheckEnough(m, userId, buildingType.Cost*float64(payload.Square)) {
		return errors.New("not enough money")
	}
	return CreateBuilding(m, userId, payload, buildingType)
}

func CreateBuilding(m *mongo.Database, userId primitive.ObjectID, payload ConstructBuildingPayload, buildingType BuildingType) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	if err := AddMoney(m, userId, (-1)*buildingType.Cost*float64(payload.Square)); err != nil {
		return err
	}

	now := time.Now()
	end := now.Add(time.Duration(float64(buildingType.BuildTime) * float64(payload.Square)))
	building := Building{
		TypeId:      payload.TypeId,
		UserId:      userId,
		X:           payload.X,
		Y:           payload.Y,
		Square:      payload.Square,
		Level:       1,
		Status:      ConstructionStatus,
		WorkStarted: now,
		WorkEnd:     end,
		HiringNeeds: 0,
		Salary:      0,
		Workers:     0,
	}

	_, err := m.Collection("buildings").InsertOne(ctx, building)
	return err
}

type FindBuildingParams struct {
	Id             *primitive.ObjectID `json:"id"`
	UserId         *primitive.ObjectID `json:"userId"`
	NickName       *string             `json:"nickName"`
	X              *int                `json:"x"`
	Y              *int                `json:"y"`
	BuildingTypeId *uint               `json:"buildingTypeId"`
	Limit          *int                `json:"limit"`
	OrderField     *string             `json:"orderField"`
	Order          *string             `json:"order"`
	Page           *int                `json:"page"`
}

type BuildingWithData struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TypeId       uint               `json:"typeId"`
	UserId       primitive.ObjectID `json:"userId"`
	X            int                `json:"x"`
	Y            int                `json:"y"`
	Square       int                `json:"square"`
	Level        int                `json:"level"`
	Status       BuildingStatus     `json:"status"`
	WorkStarted  time.Time          `json:"workStarted"`
	WorkEnd      time.Time          `json:"workEnd"`
	HiringNeeds  int                `json:"hiringNeeds"`
	Salary       float64            `json:"salary"`
	Workers      int                `json:"workers"`
	OnStrike     bool               `json:"onStrike"`
	BuildingType BuildingType       `json:"buildingType"`
	NickName     string             `json:"nickName"`
}

func GetBuildings(m *mongo.Database, findBuildingParams FindBuildingParams) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	if findBuildingParams.Id != nil {
		filter = append(filter, bson.E{Key: "buildings._id", Value: *findBuildingParams.Id})
	}
	if findBuildingParams.UserId != nil {
		filter = append(filter, bson.E{Key: "userId", Value: *findBuildingParams.UserId})
	}
	if findBuildingParams.BuildingTypeId != nil {
		filter = append(filter, bson.E{Key: "typeId", Value: *findBuildingParams.BuildingTypeId})
	}
	if findBuildingParams.X != nil {
		filter = append(filter, bson.E{Key: "x", Value: *findBuildingParams.X})
	}
	if findBuildingParams.Y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *findBuildingParams.Y})
	}

	matchStage := bson.D{{"$match", filter}}
	lookupBuildingType := bson.D{{"$lookup", bson.D{
		{"from", "buildingTypes"},
		{"localField", "typeId"},
		{"foreignField", "id"},
		{"as", "buildingType"},
	}}}
	lookupUser := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "userId"},
		{"foreignField", "_id"},
		{"as", "user"},
	}}}

	unwindBuildingType := bson.D{{"$unwind", bson.D{
		{"path", "$buildingType"},
		{"preserveNullAndEmptyArrays", true},
	}}}
	unwindUser := bson.D{{"$unwind", bson.D{
		{"path", "$user"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	projectStage := bson.D{{"$project", bson.D{
		{"typeId", 1},
		{"x", 1},
		{"y", 1},
		{"level", 1},
		{"status", 1},
		{"square", 1},
		{"buildingType.title", 1},
		{"nickName", "$user.nickName"},
	}}}

	sort := bson.D{}

	if findBuildingParams.OrderField != nil {
		if findBuildingParams.Order != nil {
			sort = append(filter, bson.E{Key: *findBuildingParams.OrderField, Value: *findBuildingParams.Order})
		} else {
			sort = append(filter, bson.E{Key: *findBuildingParams.OrderField, Value: 1})
		}
	}

	sortStage := bson.D{}

	if len(sort) != 0 {
		sortStage = bson.D{{"$sort", sort}}
	} else {
		sortStage = bson.D{{"$sort", bson.D{{"_id", -1}}}}
	}

	limit := 20
	if findBuildingParams.Limit != nil {
		limit = *findBuildingParams.Limit
	}
	limitStage := bson.D{{"$limit", limit}}

	skipStage := bson.D{{"$skip", 0}}
	if findBuildingParams.Page != nil {
		skipStage = bson.D{{"$skip", (*findBuildingParams.Page - 1) * limit}}
	}

	pipeline := mongo.Pipeline{matchStage, lookupBuildingType, lookupUser, unwindBuildingType, unwindUser,
		projectStage, sortStage, skipStage, limitStage}

	cursor, err := m.Collection("buildings").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get buildings: " + err.Error())
		return nil, err
	}

	var buildings []BuildingWithData
	if err = cursor.All(ctx, &buildings); err != nil {
		log.Println(err)
	}
	return buildings, nil
}

func GetMyBuildings(m *mongo.Database, userId primitive.ObjectID, buildingId primitive.ObjectID) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	filter = append(filter, bson.E{Key: "userId", Value: userId})
	if buildingId != primitive.NilObjectID {
		filter = append(filter, bson.E{Key: "_id", Value: buildingId})
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
		log.Println("Can't get buildings: " + err.Error())
		return nil, err
	}

	var myBuildings []BuildingWithData
	if err = cursor.All(ctx, &myBuildings); err != nil {
		log.Println(err)
	}
	return myBuildings, nil
}

func GetBuildingByID(m *mongo.Database, buildingID primitive.ObjectID) (Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var building Building
	err := m.Collection("buildings").FindOne(ctx,
		bson.M{"_id": buildingID}).Decode(&building)
	if err != nil {
		log.Println("Can't get building by Id: " + err.Error())
	}
	return building, err
}

type HiringPayload struct {
	BuildingID  primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	Salary      float64            `json:"salary" bson:"salary"`
	HiringNeeds int                `json:"hiringNeeds" bson:"hiringNeeds"`
}

func SetHiring(m *mongo.Database, userID primitive.ObjectID, payload HiringPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingByID(m, payload.BuildingID)
	if err != nil {
		return err
	}
	if userID != building.UserId && building.UserId != primitive.NilObjectID {
		return errors.New("this building doesn't belong to you")
	}
	buildingType, err := GetBuildingTypeByID(m, building.TypeId)
	if err != nil {
		return err
	}
	hiringMax := buildingType.Workers * building.Level * building.Square
	if payload.HiringNeeds > hiringMax {
		return errors.New(fmt.Sprintf("hiring needs more that maximum(%d)", hiringMax))
	}

	_, err = m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": building.Id},
		bson.M{"$set": bson.M{"salary": payload.Salary, "hiringNeeds": payload.HiringNeeds}})
	if err != nil {
		log.Println("Can't update building: " + err.Error())
	}
	return err
}

func DestroyBuilding(m *mongo.Database, userId primitive.ObjectID, buildingId primitive.ObjectID) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingByID(m, buildingId)
	if userId != building.UserId && building.UserId != primitive.NilObjectID {
		return errors.New("for attempting to destroy someone else's building, inevitable punishment awaits you")
	}
	if err != nil {
		log.Println("Can't destroy building: " + err.Error())
		return err
	}

	_, err = m.Collection("buildings").DeleteOne(ctx, bson.M{"_id": buildingId, "userId": userId})
	if err != nil {
		log.Println("Failed to delete building: " + err.Error())
		return err
	}

	return nil
}

func GetAllReadyStorages(m *mongo.Database) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var readyStorages []Building

	filter := bson.M{"status": ReadyStatus, "onStrike": false, "typeId": 1}
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return readyStorages, err
	}

	err = cursor.All(ctx, &readyStorages)
	return readyStorages, err
}

func BuildingStatusUpdate(m *mongo.Database, buildingId primitive.ObjectID, status BuildingStatus) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": buildingId},
		bson.M{"$set": bson.M{"status": status}})
	return err
}

func GetBuildingsForHiring(m *mongo.Database) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.M{"salary": bson.M{"$ne": 0}, "hiringNeeds": bson.M{"$ne": 0}}
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var buildings []Building
	err = cursor.All(ctx, &buildings)
	return buildings, err
}

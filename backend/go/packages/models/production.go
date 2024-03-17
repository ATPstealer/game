package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Production struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingId  primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	BlueprintId uint               `json:"blueprintId" bson:"blueprintId"`
	WorkStarted time.Time          `json:"workStarted" bson:"workStarted"`
	WorkEnd     time.Time          `json:"workEnd" bson:"workEnd"`
}

type StartWorkPayload struct {
	BuildingId  primitive.ObjectID
	BlueprintId uint
	Duration    time.Duration
}

func StartWork(m *mongo.Database, userId primitive.ObjectID, payload StartWorkPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		log.Println("Can't find buildings: " + err.Error())
		return err
	}
	if building.Status != ReadyStatus {
		return errors.New("Building not ready. Status is " + string(building.Status))
	}
	if building.UserId != userId {
		err := errors.New("this building don't belong you")
		log.Println(err)
		return err
	}
	blueprintResult, err := GetBlueprintById(m, payload.BlueprintId)
	if err != nil {
		log.Println("invalid blueprint" + err.Error())
		return err
	}
	if blueprintResult.ProducedInId != building.TypeId {
		err := errors.New("can't product it here")
		log.Println(err)
		return err
	}
	log.Println(building.WorkStarted)
	now := time.Now()
	end := now.Add(payload.Duration)

	prod := Prod{
		BlueprintId: payload.BlueprintId,
	}

	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, bson.M{
		"$set": bson.M{
			"status":      ProductionStatus,
			"workStarted": now,
			"workEnd":     end,
			"prod":        &prod,
		},
	})
	if err != nil {
		log.Println("Failed to update building: " + err.Error())
		return err
	}

	_, err = m.Collection("productions").InsertOne(ctx, Production{
		BuildingId:  payload.BuildingId,
		BlueprintId: payload.BlueprintId,
		WorkStarted: now,
		WorkEnd:     end,
	})
	if err != nil {
		log.Println("Failed to insert production: " + err.Error())
		return err
	}
	return nil
}

func ProductionSetWorkStarted(m *mongo.Database, productionId primitive.ObjectID, timeStart time.Time) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("productions").UpdateOne(ctx,
		bson.M{"_id": productionId},
		bson.M{"$set": bson.M{"workStarted": timeStart}})
	return err
}

type ProductionWithData struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingId   primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	BlueprintId  uint               `json:"blueprintId" bson:"blueprintId"`
	WorkStarted  time.Time          `json:"workStarted" bson:"workStarted"`
	WorkEnd      time.Time          `json:"workEnd" bson:"workEnd"`
	Building     Building           `json:"building" bson:"building"`
	BuildingType BuildingType       `json:"buildingType" bson:"buildingType"`
}

func GetProduction(m *mongo.Database) ([]ProductionWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"workEnd", bson.D{{"$gt", time.Now()}}}}
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
	cursor, err := m.Collection("productions").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get productions: " + err.Error())
		return nil, err
	}

	var productions []ProductionWithData
	if err = cursor.All(ctx, &productions); err != nil {
		log.Println(err)
	}
	return productions, nil
}

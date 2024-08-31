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
	BlueprintId uint `json:"blueprintId" bson:"blueprintId"`
}

type StartProductionPayload struct {
	BuildingId  primitive.ObjectID
	BlueprintId uint
	Duration    time.Duration
}

func GetBuildingsProduction(m *mongo.Database) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{
		{"workEnd", bson.D{{"$gt", time.Now()}}},
		{"production", bson.D{{"$ne", nil}}},
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

func StartProduction(m *mongo.Database, userId primitive.ObjectID, payload StartProductionPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		log.Println("Can't find buildings: " + err.Error())
		return err
	}
	if building.Status != ReadyStatus {
		return errors.New("building busy")
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
		return err
	}

	now := time.Now()
	end := now.Add(payload.Duration)

	production := Production{
		BlueprintId: payload.BlueprintId,
	}

	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, bson.M{
		"$set": bson.M{
			"status":      ProductionStatus,
			"workStarted": now,
			"workEnd":     end,
			"production":  &production,
		},
	})
	if err != nil {
		log.Println("Failed to update building: " + err.Error())
		return err
	}
	return nil
}

func StopProduction(m *mongo.Database, userId primitive.ObjectID, payload StartProductionPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		log.Println("Can't find buildings: " + err.Error())
		return err
	}

	if building.UserId != userId {
		err := errors.New("this building don't belong you")
		log.Println(err)
		return err
	}

	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, bson.M{
		"$set": bson.M{
			"status":      ReadyStatus,
			"production":  nil,
			"workStarted": nil,
			"workEnd":     nil,
		},
	})

	return err
}

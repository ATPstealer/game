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

type ProductionMongo struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingID  primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	BlueprintID uint               `json:"blueprintId" bson:"blueprintId"`
	WorkStarted *time.Time         `json:"workStarted" bson:"workStarted"`
	WorkEnd     *time.Time         `json:"workEnd" bson:"workEnd"`
}

type StartWorkPayloadMongo struct {
	BuildingID  primitive.ObjectID
	BlueprintID uint
	Duration    time.Duration
}

func StartWorkMongo(m *mongo.Database, userID primitive.ObjectID, payload StartWorkPayloadMongo) error {
	building, err := GetBuildingByIDMongo(m, payload.BuildingID)
	if err != nil {
		log.Println("Can't find buildings: " + err.Error())
		return err
	}
	if building.Status != ReadyStatus {
		return errors.New("Building not ready. Status is " + string(building.Status))
	}
	if building.UserID != userID {
		err := errors.New("this building don't belong you")
		log.Println(err)
		return err
	}
	blueprintResult, err := GetBlueprintByIDMongo(m, payload.BlueprintID)
	if err != nil {
		log.Println("invalid blueprint" + err.Error())
		return err
	}
	if blueprintResult.ProducedInID != building.TypeID {
		err := errors.New("can't product it here")
		log.Println(err)
		return err
	}
	log.Println(building.WorkStarted)
	now := time.Now()
	end := now.Add(payload.Duration)

	_, err = m.Collection("buildings").UpdateOne(context.TODO(), bson.M{"_id": building.ID}, bson.M{
		"$set": bson.M{
			"status":      ProductionStatus,
			"workStarted": &now,
			"workEnd":     &end,
		},
	})
	if err != nil {
		log.Println("Failed to update building: " + err.Error())
		return err
	}

	_, err = m.Collection("productions").InsertOne(context.TODO(), ProductionMongo{
		BuildingID:  payload.BuildingID,
		BlueprintID: payload.BlueprintID,
		WorkStarted: &now,
		WorkEnd:     &end,
	})
	if err != nil {
		log.Println("Failed to insert production: " + err.Error())
		return err
	}
	return nil
}

func ProductionSetWorkStarted(m *mongo.Database, productionId primitive.ObjectID, time *time.Time) error {
	_, err := m.Collection("productions").UpdateOne(context.TODO(),
		bson.M{"_id": productionId},
		bson.M{"$set": bson.M{"workStarted": &time}})
	return err
}

type ProductionWithDataMongo struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BuildingID   primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	BlueprintID  uint               `json:"blueprintId" bson:"blueprintId"`
	WorkStarted  *time.Time         `json:"workStarted" bson:"workStarted"`
	WorkEnd      *time.Time         `json:"workEnd" bson:"workEnd"`
	Building     BuildingMongo      `json:"building" bson:"building"`
	BuildingType BuildingTypeMongo  `json:"buildingType" bson:"buildingType"`
}

func GetProductionMongo(m *mongo.Database) ([]ProductionWithDataMongo, error) {
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
	cursor, err := m.Collection("productions").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get productions: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var productions []ProductionWithDataMongo
	if err = cursor.All(context.TODO(), &productions); err != nil {
		log.Println(err)
	}
	return productions, nil
}

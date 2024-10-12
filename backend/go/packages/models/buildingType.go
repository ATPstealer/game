package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type BuildingType struct {
	Id               uint          `bson:"id" json:"id"`
	Title            string        `bson:"title" json:"title"`
	Description      string        `bson:"description" json:"description"`
	Cost             float64       `bson:"cost" json:"cost"`
	Requirements     string        `bson:"requirements" json:"requirements"`
	BuildTime        time.Duration `bson:"buildTime" json:"buildTime"`
	BuildingGroup    string        `bson:"buildingGroup" json:"buildingGroup"`
	BuildingSubGroup string        `bson:"buildingSubGroup" json:"buildingSubGroup"`
	Capacity         float64       `bson:"capacity" json:"capacity"`
	Workers          int           `bson:"workers" json:"workers"`
} // @name buildingType

func GetAllBuildingTypes(m *mongo.Database) ([]BuildingType, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var buildingTypes []BuildingType
	cursor, err := m.Collection("buildingTypes").Find(ctx, bson.M{})
	if err != nil {
		return buildingTypes, err
	}

	err = cursor.All(ctx, &buildingTypes)
	return buildingTypes, err
}

func GetBuildingTypeById(m *mongo.Database, typeId uint) (BuildingType, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var buildingType BuildingType
	res := m.Collection("buildingTypes").FindOne(ctx, bson.M{"id": typeId})
	err := res.Decode(&buildingType)
	return buildingType, err
}

func GetBuildingTypesByBuildingGroup(m *mongo.Database, group string) ([]BuildingType, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var buildingTypes []BuildingType
	cursor, err := m.Collection("buildingTypes").Find(ctx, bson.M{"buildingGroup": group})
	if err != nil {
		return buildingTypes, err
	}

	err = cursor.All(ctx, &buildingTypes)
	return buildingTypes, err
}

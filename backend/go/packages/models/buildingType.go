package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type BuildingType struct {
	Id               uint             `bson:"id" json:"id" validate:"required"`
	Title            string           `bson:"title" json:"title" validate:"required"`
	Description      string           `bson:"description" json:"description" validate:"required"`
	Cost             float64          `bson:"cost" json:"cost" validate:"required"`
	Requirements     []ResourceAmount `bson:"requirements" json:"requirements" validate:"required"`
	BuildTime        time.Duration    `bson:"buildTime" json:"buildTime" validate:"required"`
	BuildingGroup    string           `bson:"buildingGroup" json:"buildingGroup" validate:"required"`
	BuildingSubGroup string           `bson:"buildingSubGroup" json:"buildingSubGroup" validate:"required"`
	Capacity         float64          `bson:"capacity" json:"capacity" validate:"required"`
	Workers          int              `bson:"workers" json:"workers" validate:"required"`
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

package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type BuildingTypeMongo struct {
	ID               uint          `bson:"id" json:"id"`
	Title            string        `bson:"title" json:"title"`
	Description      string        `bson:"description" json:"description"`
	Cost             float64       `bson:"cost" json:"cost"`
	Requirements     string        `bson:"requirements" json:"requirements"`
	BuildTime        time.Duration `bson:"buildTime" json:"buildTime"`
	BuildingGroup    string        `bson:"buildingGroup" json:"buildingGroup"`
	BuildingSubGroup string        `bson:"buildingSubGroup" json:"buildingSubGroup"`
	Capacity         float64       `bson:"capacity" json:"capacity"`
	Workers          int           `bson:"workers" json:"workers"`
}

func GetAllBuildingTypesMongo(m *mongo.Database) ([]BuildingTypeMongo, error) {
	var buildingTypes []BuildingTypeMongo
	cursor, err := m.Collection("buildingTypes").Find(context.Background(), bson.M{})
	if err != nil {
		return buildingTypes, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.TODO(), &buildingTypes)
	return buildingTypes, err
}

func GetBuildingTypeByIDMongo(m *mongo.Database, typeID uint) (BuildingTypeMongo, error) {
	var buildingType BuildingTypeMongo
	res := m.Collection("buildingTypes").FindOne(context.Background(), bson.M{"id": typeID})
	err := res.Decode(&buildingType)
	return buildingType, err
}

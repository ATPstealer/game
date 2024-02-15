package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"time"
)

type BuildingType struct {
	gorm.Model                     // This includes some common fields like ID, CreatedAt, UpdatedAt, and DeletedAt.
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	Cost             float64       `json:"cost"`
	Requirements     string        `json:"requirements"` // TODO: Resource object
	BuildTime        time.Duration `json:"buildTime"`
	BuildingGroup    string        `json:"buildingGroup"`
	BuildingSubGroup string        `json:"buildingSubGroup"`
	Capacity         float64       `json:"capacity"`
	Workers          int           `json:"workers"`
}

type BuildingTypeResult struct {
	ID               uint          `json:"id"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	Cost             float64       `json:"cost"`
	Requirements     string        `json:"requirements"`
	BuildTime        time.Duration `json:"buildTime"`
	BuildingGroup    string        `json:"buildingGroup"`
	BuildingSubGroup string        `json:"buildingSubGroup"`
	Capacity         float64       `json:"capacity"`
	Workers          int           `json:"workers"`
}

func GetBuildingTypeByID(db *gorm.DB, typeID uint) (BuildingTypeResult, error) {
	var buildingType BuildingTypeResult
	res := db.Model(&BuildingType{}).Where("id = ?", typeID).First(&buildingType)
	if res.Error != nil {
		log.Println("Can't get Building Type: " + res.Error.Error())
	}
	return buildingType, res.Error
}

func GetAllBuildingTypes(db *gorm.DB) ([]BuildingTypeResult, error) {
	var buildingTypes []BuildingTypeResult
	res := db.Model(&BuildingType{}).Find(&buildingTypes).Select("title")
	if res.Error != nil {
		log.Println("Can't get Building Type: " + res.Error.Error())
	}
	return buildingTypes, res.Error
}

// mongo

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

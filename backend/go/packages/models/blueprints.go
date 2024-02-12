package models

import (
	"context"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"time"
)

type Blueprint struct {
	gorm.Model
	Name              string        `json:"name"`
	ProducedResources string        `gorm:"type:json" json:"producedResources"`
	UsedResources     string        `gorm:"type:json" json:"usedResources"`
	ProducedInID      uint          `json:"producedInId"`
	ProductionTime    time.Duration `json:"productionTime"`
}

type ResourceAmount struct {
	ResourceID uint    `json:"resourceId"`
	Amount     float64 `json:"amount"`
}

type BlueprintResult struct {
	ID                uint             `json:"id"`
	Name              string           `json:"name"`
	ProducedResources []ResourceAmount `json:"producedResources"`
	UsedResources     []ResourceAmount `json:"usedResources"`
	ProducedInID      uint             `json:"producedInId"`
	ProductionTime    time.Duration    `json:"productionTime"`
}

func GetBlueprintByID(db *gorm.DB, blueprintID uint) (BlueprintResult, error) {
	var blueprint Blueprint
	res := db.Model(&Blueprint{}).Where("id = ?", blueprintID).First(&blueprint)
	if res.Error != nil {
		log.Println("Can't get Blueprint: " + res.Error.Error())
	}
	blueprintResult, err := BlueprintToBlueprintResult(blueprint)
	if err != nil {
		log.Println(err.Error())
	}
	return blueprintResult, err
}

func GetBlueprints(db *gorm.DB, blueprintID uint) ([]BlueprintResult, error) {
	var blueprints []Blueprint
	if blueprintID != 0 {
		db.Where("id = ?", blueprintID).Find(&blueprints)
	} else {
		db.Find(&blueprints)
	}
	var blueprintResultArray []BlueprintResult
	for _, blueprint := range blueprints {
		blueprintResult, err := BlueprintToBlueprintResult(blueprint)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		blueprintResultArray = append(blueprintResultArray, blueprintResult)
	}
	return blueprintResultArray, nil
}

func BlueprintToBlueprintResult(blueprint Blueprint) (blueprintResult BlueprintResult, err error) {
	var producedResources []ResourceAmount
	if err := json.Unmarshal([]byte(blueprint.ProducedResources), &producedResources); err != nil {
		log.Println("Can't make JSON from ProducedResources field")
		return BlueprintResult{}, err
	}
	var usedResources []ResourceAmount
	if err := json.Unmarshal([]byte(blueprint.UsedResources), &usedResources); err != nil {
		log.Println("Can't make JSON from ProducedResources field")
		return BlueprintResult{}, err
	}
	return BlueprintResult{
		ID:                blueprint.ID,
		Name:              blueprint.Name,
		ProducedResources: producedResources,
		UsedResources:     usedResources,
		ProducedInID:      blueprint.ProducedInID,
		ProductionTime:    blueprint.ProductionTime,
	}, nil
}

//mongo

type BlueprintMongo struct {
	ID                uint                  `json:"id" bson:"id"`
	Name              string                `json:"name" bson:"name"`
	ProducedResources []ResourceAmountMongo `json:"producedResources" bson:"producedResources"`
	UsedResources     []ResourceAmountMongo `json:"usedResources" bson:"usedResources"`
	ProducedInID      uint                  `json:"producedInId" bson:"producedInId"`
	ProductionTime    time.Duration         `json:"productionTime" bson:"productionTime"`
}

type ResourceAmountMongo struct {
	ResourceID uint    `json:"resourceId" bson:"resourceId"`
	Amount     float64 `json:"amount" bson:"amount"`
}

func GetBlueprintsMongo(m *mongo.Database, blueprintID uint) ([]BlueprintMongo, error) {
	var blueprints []BlueprintMongo
	filter := bson.M{}
	if blueprintID != 0 {
		filter["id"] = blueprintID
	}
	cur, err := m.Collection("blueprints").Find(context.TODO(), filter)
	if err != nil {
		log.Println("Can't get blueprints: " + err.Error())
		return nil, err
	}
	defer cur.Close(context.TODO())

	err = cur.All(context.TODO(), &blueprints)
	return blueprints, err
}

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

type Logistics struct {
	CapacityMax float64 `json:"capacityMax" bson:"capacityMax"`
	Capacity    float64 `json:"capacity" bson:"capacity"`
	Speed       float64 `json:"speed" bson:"speed"`
	Price       float64 `json:"price" bson:"price"`
	Revenue     float64 `json:"revenue" bson:"revenue"`
}

// TODO: переделать все поиска зданий на этот вариант
func GetAllReadyLogisticsHubs(m *mongo.Database) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var readyLogisticsHubs []Building

	logisticsTypes, err := GetBuildingTypesByBuildingGroup(m, "Logistics")
	if err != nil {
		return readyLogisticsHubs, err
	}

	var typeIds []uint
	for _, logisticsType := range logisticsTypes {
		typeIds = append(typeIds, logisticsType.Id)
	}

	filter := bson.M{"status": ReadyStatus, "onStrike": false, "typeId": bson.M{"$in": typeIds}}
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return readyLogisticsHubs, err
	}

	err = cursor.All(ctx, &readyLogisticsHubs)
	return readyLogisticsHubs, err
}

func LogisticsReset(m *mongo.Database, building Building) error {
	logisticsEquipmentEffect, err := findLogisticsEffect(m, building.EquipmentEffect)
	if err != nil {
		return err
	}

	price := 0.0
	if building.Logistics != nil {
		price = building.Logistics.Price
	}

	logistics := Logistics{
		CapacityMax: logisticsEquipmentEffect.Value,
		Capacity:    logisticsEquipmentEffect.Value,
		Speed:       logisticsEquipmentEffect.ValueSecond,
		Price:       price,
		Revenue:     0.0,
	}
	err = updateBuildingLogistics(m, building.Id, logistics)

	return err
}

func findLogisticsEffect(m *mongo.Database, equipmentEffects *[]EquipmentEffect) (EquipmentEffect, error) {
	if equipmentEffects != nil {
		for index, equipmentEffect := range *equipmentEffects {
			if equipmentEffect.EffectId == LogisticsCapacity {
				return (*equipmentEffects)[index], nil
			}
		}
	}
	return EquipmentEffect{}, errors.New("logistics effect not found")
}

func updateBuildingLogistics(m *mongo.Database, buildingId primitive.ObjectID, logistics Logistics) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var err error
	if logistics.CapacityMax != 0 {
		_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": buildingId}, bson.M{
			"$set": bson.M{
				"logistics": logistics,
			},
		})
	} else {
		_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": buildingId}, bson.M{
			"$unset": bson.M{
				"logistics": "",
			},
		})
	}

	if err != nil {
		log.Println("Failed to update building: " + err.Error())
		return err
	}
	return nil
}

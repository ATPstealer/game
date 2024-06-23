package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// TODO добавить учет количества рабочих и количество оборудования.

type EquipmentEffect struct {
	EffectId    uint    `json:"effectId" bson:"effectId"`
	BlueprintId uint    `json:"blueprintId" bson:"blueprintId"`
	Value       float64 `json:"value" bson:"value"`
}

func countEffects(m *mongo.Database, buildingId primitive.ObjectID) error {
	building, err := GetBuildingById(m, buildingId)
	if err != nil {
		return err
	}

	equipmentEffects := new([]EquipmentEffect)
	if building.Equipment != nil {
		for _, equipment := range *building.Equipment {
			log.Println(equipment)
			equipmentType, err := GetEquipmentTypesByID(m, equipment.EquipmentTypeId)
			if err != nil {
				return err
			}

			if len(equipmentType.BlueprintIds) == 0 {
				addEffect(equipmentEffects, equipmentType.EffectId, 0, equipmentType.Value*float64(equipment.Amount))
			} else {
				for _, blueprintId := range equipmentType.BlueprintIds {
					addEffect(equipmentEffects, equipmentType.EffectId, blueprintId, equipmentType.Value*float64(equipment.Amount))
				}
			}
		}
	} else {
		equipmentEffects = nil
	}

	return saveEquipmentEffects(m, equipmentEffects, buildingId)
}

func addEffect(equipmentEffects *[]EquipmentEffect, effectId uint, blueprintId uint, value float64) {
	for i, effect := range *equipmentEffects {
		if effect.EffectId == effectId && effect.BlueprintId == blueprintId {
			(*equipmentEffects)[i].Value += value
			return
		}
	}
	*equipmentEffects = append(*equipmentEffects, EquipmentEffect{
		EffectId:    effectId,
		BlueprintId: blueprintId,
		Value:       value,
	})
}

func saveEquipmentEffects(m *mongo.Database, equipmentEffects *[]EquipmentEffect, buildingId primitive.ObjectID) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": buildingId}, bson.M{"$set": bson.M{"equipmentEffect": equipmentEffects}})
	if err != nil {
		log.Println("Error updating building with equipment effects:", err)
	}
	return err
}

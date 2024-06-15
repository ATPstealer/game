package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type EquipmentType struct {
	Id             uint    `json:"id" bson:"id"`
	Name           string  `json:"name" bson:"name"`
	ResourceTypeId uint    `json:"resourceTypeId" bson:"resourceTypeId"`
	Durability     int     `json:"durability" bson:"durability"`
	BlueprintIds   []uint  `json:"blueprintIds" bson:"blueprintIds"`
	EffectId       uint    `json:"effectId" bson:"effectId"`
	Value          float64 `json:"value" bson:"value"`
	Square         float64 `json:"square" bson:"square"`
}

func GetAllEquipmentTypes(m *mongo.Database) ([]EquipmentType, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var equipmentTypes []EquipmentType
	cursor, err := m.Collection("equipmentTypes").Find(ctx, bson.M{})
	if err != nil {
		return equipmentTypes, err
	}

	err = cursor.All(ctx, &equipmentTypes)
	return equipmentTypes, err
}

func GetEquipmentTypesByID(m *mongo.Database, typeId uint) (EquipmentType, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var equipmentType EquipmentType
	err := m.Collection("EquipmentTypes").FindOne(ctx, bson.M{"id": typeId}).Decode(&equipmentType)
	if err != nil {
		log.Println("Can't get resource type: " + err.Error())
	}
	return equipmentType, err
}

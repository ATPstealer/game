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

type FindLogisticsParams struct {
	X           *int
	Y           *int
	MinCapacity *float64
}

type LogisticsWithData struct {
	X          int                `json:"x"`
	Y          int                `json:"y"`
	BuildingId primitive.ObjectID `json:"buildingId"`
	Capacity   float64            `json:"capacity"`
	Speed      float64            `json:"speed"`
	Price      float64            `json:"price"`
}

func GetLogisticsCapacity(m *mongo.Database, findLogisticsParams FindLogisticsParams) ([]LogisticsWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	if findLogisticsParams.X != nil {
		filter = append(filter, bson.E{Key: "x", Value: *findLogisticsParams.X})
	}
	if findLogisticsParams.Y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *findLogisticsParams.Y})
	}
	if findLogisticsParams.MinCapacity != nil {
		filter = append(filter, bson.E{Key: "logistics.capacity", Value: bson.D{{Key: "$gte", Value: *findLogisticsParams.MinCapacity}}})
	}

	var logistics []LogisticsWithData
	var buildings []Building
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return logistics, err
	}
	if err = cursor.All(ctx, &buildings); err != nil {
		return logistics, err
	}

	logistics = getLogistics(buildings)
	return logistics, err
}

func getLogistics(buildings []Building) []LogisticsWithData {
	var logistics []LogisticsWithData
	for _, building := range buildings {
		if building.Logistics != nil {
			logistics = append(logistics, LogisticsWithData{
				X:          building.X,
				Y:          building.Y,
				BuildingId: building.Id,
				Capacity:   building.Logistics.Capacity,
				Speed:      building.Logistics.Speed,
				Price:      building.Logistics.Price,
			})
		}
	}
	return logistics
}

func GetLogisticsPriceAndSpeed(m *mongo.Database, buildingId primitive.ObjectID) (float64, float64) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var building Building
	err := m.Collection("buildings").FindOne(ctx, bson.M{"_id": buildingId}).Decode(&building)
	if err != nil {
		log.Println("Can't get building by Id: " + err.Error())
		return 0.0, 0.0
	}

	if building.Logistics != nil {
		return building.Logistics.Price, building.Logistics.Speed
	}
	return 0.0, 0.0
}

func CheckEnoughCapacity(m *mongo.Database, buildingId primitive.ObjectID, capacity float64) bool {
	building, err := GetBuildingById(m, buildingId)
	if err != nil {
		log.Println("Can't get building by Id: " + err.Error())
		return false
	}
	return building.Logistics.Capacity >= capacity
}

func WithdrawLogisticsCapacity(m *mongo.Database, buildingId primitive.ObjectID, capacity float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, buildingId)
	if err != nil {
		log.Println("Can't get building by Id: " + err.Error())
		return err
	}
	if building.Logistics == nil {
		return errors.New("not enough capacity in this hub")
	}
	if building.Logistics.Capacity < capacity {
		return errors.New("not enough capacity in this hub")
	}

	building.Logistics.Capacity -= capacity
	building.Logistics.Revenue += capacity * building.Logistics.Price
	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, bson.M{
		"$set": bson.M{
			"logistics": &building.Logistics,
		},
	})
	if err != nil {
		log.Println("Failed to update building: " + err.Error())
		return err
	}
	return nil
}

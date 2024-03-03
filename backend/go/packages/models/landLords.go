package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type LandLordMongo struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Square int                `json:"square"`
	X      int                `json:"x"`
	Y      int                `json:"y"`
}

type BuyLandPayload struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Square int `json:"square"`
}

func BuyLandMongo(m *mongo.Database, userID primitive.ObjectID, payload BuyLandPayload) (float64, error) {
	occupiedLand, err := GetCellOccupiedLandMongo(m, payload.X, payload.Y)
	if err != nil {
		return 0, err
	}
	price := 10 * (float64(occupiedLand)*2 + 1 + float64(payload.Square)) * float64(payload.Square) / 2

	if !CheckEnoughMoneyMongo(m, userID, price) {
		return 0, errors.New("not enough money")
	}
	enoughLand, err := CheckEnoughLandMongo(m, payload.X, payload.Y, payload.Square)
	if err != nil {
		return 0, err
	}
	if !enoughLand {
		return 0, errors.New("not enough land")
	}

	if err := AddMoneyMongo(m, userID, (-1)*price); err != nil {
		return 0, err
	}

	if err := AddCivilSavingsMongo(m, payload.X, payload.Y, price); err != nil {
		log.Println("Can't add civil money" + err.Error())
	}

	_, err = m.Collection("landLords").UpdateOne(context.TODO(),
		bson.M{
			"userId": userID,
			"x":      payload.X,
			"y":      payload.Y,
		},
		bson.M{
			"$inc": bson.M{
				"square": payload.Square,
			},
			"$setOnInsert": bson.M{
				"userId": userID,
				"x":      payload.X,
				"y":      payload.Y,
			},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return 0, err
	}
	return price, nil
}

func GetCellOwnersMongo(m *mongo.Database, x int, y int) ([]LandLordMongo, error) {
	cursor, err := m.Collection("landLords").Find(context.TODO(), bson.M{"x": x, "y": y})
	if err != nil {
		return nil, fmt.Errorf("failed to execute mongoDB query: %w", err)
	}
	defer cursor.Close(context.Background())

	var landLords []LandLordMongo
	if err = cursor.All(context.TODO(), &landLords); err != nil {
		return nil, err
	}

	return landLords, nil
}

func GetAllLandLordsMongo(m *mongo.Database) ([]LandLordMongo, error) {
	cursor, err := m.Collection("landLords").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to execute mongoDB query: %w", err)
	}
	defer cursor.Close(context.Background())

	var landLords []LandLordMongo
	if err = cursor.All(context.TODO(), &landLords); err != nil {
		return nil, err
	}

	return landLords, nil
}

func GetMyLandsMongo(m *mongo.Database, userID primitive.ObjectID) ([]LandLordMongo, error) {
	cursor, err := m.Collection("landLords").Find(context.TODO(), bson.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to execute mongoDB query: %w", err)
	}
	defer cursor.Close(context.Background())

	var landLords []LandLordMongo
	if err = cursor.All(context.TODO(), &landLords); err != nil {
		return nil, err
	}

	return landLords, nil
}

func CheckEnoughLandForBuildingMongo(m *mongo.Database, userID primitive.ObjectID, square int, x int, y int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var myLandInCell LandLordMongo
	err := m.Collection("landLords").FindOne(ctx,
		bson.M{
			"userId": userID,
			"x":      x,
			"y":      y,
		}).Decode(&myLandInCell)
	if err != nil {
		log.Println("Can't get my cell lands: " + err.Error())
		return false, err
	}

	var myBuildingsInCell []BuildingMongo

	cursor, err := m.Collection("buildings").Find(ctx,
		bson.M{
			"userId": userID,
			"x":      x,
			"y":      y,
		})
	if err != nil {
		log.Println("Can't get Buildings in Cell: " + err.Error())
		return false, err
	}

	if err = cursor.All(context.TODO(), &myBuildingsInCell); err != nil {
		return false, err
	}

	cursor.Close(ctx)

	freeLand := myLandInCell.Square
	for _, building := range myBuildingsInCell {
		freeLand -= building.Square
	}

	return freeLand >= square, nil
}

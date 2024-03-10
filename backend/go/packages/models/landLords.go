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

type LandLord struct {
	Id     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Square int                `json:"square"`
	X      int                `json:"x"`
	Y      int                `json:"y"`
}

type BuyLandPayload struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Square int `json:"square"`
}

func BuyLand(m *mongo.Database, userId primitive.ObjectID, payload BuyLandPayload) (float64, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	occupiedLand, err := GetCellOccupiedLand(m, payload.X, payload.Y)
	if err != nil {
		return 0, err
	}
	price := 10 * (float64(occupiedLand)*2 + 1 + float64(payload.Square)) * float64(payload.Square) / 2

	if !CheckEnough(m, userId, price) {
		return 0, errors.New("not enough money")
	}
	enoughLand, err := CheckEnoughLand(m, payload.X, payload.Y, payload.Square)
	if err != nil {
		return 0, err
	}
	if !enoughLand {
		return 0, errors.New("not enough land")
	}

	if err := AddMoney(m, userId, (-1)*price); err != nil {
		return 0, err
	}

	if err := AddCivilSavings(m, payload.X, payload.Y, price); err != nil {
		log.Println("Can't add civil money" + err.Error())
	}

	_, err = m.Collection("landLords").UpdateOne(ctx,
		bson.M{
			"userId": userId,
			"x":      payload.X,
			"y":      payload.Y,
		},
		bson.M{
			"$inc": bson.M{
				"square": payload.Square,
			},
			"$setOnInsert": bson.M{
				"userId": userId,
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

func GetCellOwners(m *mongo.Database, x int, y int) ([]LandLord, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	cursor, err := m.Collection("landLords").Find(ctx, bson.M{"x": x, "y": y})
	if err != nil {
		return nil, fmt.Errorf("failed to execute mongoDB query: %w", err)
	}

	var landLords []LandLord
	if err = cursor.All(ctx, &landLords); err != nil {
		return nil, err
	}

	return landLords, nil
}

func GetAllLandLords(m *mongo.Database) ([]LandLord, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	cursor, err := m.Collection("landLords").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to execute mongoDB query: %w", err)
	}

	var landLords []LandLord
	if err = cursor.All(ctx, &landLords); err != nil {
		return nil, err
	}

	return landLords, nil
}

func GetMyLands(m *mongo.Database, userId primitive.ObjectID) ([]LandLord, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	cursor, err := m.Collection("landLords").Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return nil, fmt.Errorf("failed to execute mongoDB query: %w", err)
	}

	var landLords []LandLord
	if err = cursor.All(ctx, &landLords); err != nil {
		return nil, err
	}

	return landLords, nil
}

func CheckEnoughLandForBuilding(m *mongo.Database, userId primitive.ObjectID, square int, x int, y int) (bool, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var myLandInCell LandLord
	err := m.Collection("landLords").FindOne(ctx,
		bson.M{
			"userId": userId,
			"x":      x,
			"y":      y,
		}).Decode(&myLandInCell)
	if err != nil {
		log.Println("Can't get my cell lands: " + err.Error())
		return false, err
	}

	var myBuildingsInCell []Building

	cursor, err := m.Collection("buildings").Find(ctx,
		bson.M{
			"userId": userId,
			"x":      x,
			"y":      y,
		})
	if err != nil {
		log.Println("Can't get Buildings in Cell: " + err.Error())
		return false, err
	}

	if err = cursor.All(ctx, &myBuildingsInCell); err != nil {
		return false, err
	}

	freeLand := myLandInCell.Square
	for _, building := range myBuildingsInCell {
		freeLand -= building.Square
	}

	return freeLand >= square, nil
}

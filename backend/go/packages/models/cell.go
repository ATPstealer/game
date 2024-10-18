package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Cell struct {
	Id               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	CellName         string             `json:"cellName" bson:"cellName" validate:"required"`
	X                int                `json:"x" bson:"x" validate:"required"`
	Y                int                `json:"y" bson:"y" validate:"required"`
	SurfaceImagePath string             `json:"surfaceImagePath" bson:"surfaceImagePath" validate:"required"`
	Square           int                `json:"square" bson:"square" validate:"required"`
	Pollution        float64            `json:"pollution" bson:"pollution" validate:"required"`
	Population       float64            `json:"population" bson:"population" validate:"required"`
	CivilSavings     float64            `json:"civilSavings" bson:"civilSavings validate:"required"`
	SpendRate        float64            `json:"SpendRate" bson:"spendRate" validate:"required"`
	Education        float64            `json:"education" bson:"education" validate:"required"`
	Crime            float64            `json:"crime" bson:"crime" validate:"required"`
	Medicine         float64            `json:"medicine" bson:"medicine" validate:"required"`
	AverageSalary    float64            `json:"averageSalary" bson:"averageSalary" validate:"required"`
} // @name cell

func CheckEnoughLand(m *mongo.Database, x int, y int, squareForBuy int) (bool, error) {
	if squareForBuy <= 0 {
		return false, errors.New("square should be greater than 0")
	}

	occupiedLand, err := GetCellOccupiedLand(m, x, y)
	if err != nil {
		return false, err
	}

	cell, err := GetCell(m, x, y)
	if err != nil {
		return false, err
	}

	return cell.Square-occupiedLand >= squareForBuy, nil
}

func GetCell(m *mongo.Database, x int, y int) (Cell, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var cell Cell
	err := m.Collection("cells").FindOne(ctx,
		bson.M{"x": x, "y": y}).Decode(&cell)
	return cell, err
}

func GetCellOccupiedLand(m *mongo.Database, x int, y int) (int, error) {
	landLords, err := GetCellOwners(m, x, y)
	if err != nil {
		return 0, err
	}

	occupiedLand := 0
	for _, landLord := range landLords {
		occupiedLand += landLord.Square
	}
	return occupiedLand, nil
}

func AddCivilSavings(m *mongo.Database, x int, y int, money float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("cells").UpdateOne(ctx,
		bson.M{
			"x": x,
			"y": y,
		},
		bson.M{
			"$inc": bson.M{
				"civilSavings": money,
			},
		})
	return err
}

func GetAllCells(m *mongo.Database) ([]Cell, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var cells []Cell

	cursor, err := m.Collection("cells").Find(ctx, bson.M{})
	if err != nil {
		return cells, err
	}

	err = cursor.All(ctx, &cells)
	return cells, err
}

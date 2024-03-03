package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Cell struct {
	CellName         string  `bson:"cellName" json:"cellName"`
	X                int     `bson:"x" json:"x"`
	Y                int     `bson:"y" json:"y"`
	SurfaceImagePath string  `bson:"surfaceImagePath" json:"surfaceImagePath"`
	Square           int     `bson:"square" json:"square"`
	Pollution        float64 `bson:"pollution" json:"pollution"`
	Population       float64 `bson:"population" json:"population"`
	CivilSavings     float64 `bson:"civilSavings" json:"civilSavings"`
	SpendRate        float64 `bson:"spendRate" json:"SpendRate"`
	Education        float64 `bson:"education" json:"education"`
	Crime            float64 `bson:"crime" json:"crime"`
	Medicine         float64 `bson:"medicine" json:"medicine"`
	AverageSalary    float64 `bson:"averageSalary" json:"averageSalary"`
}

func CheckEnoughLand(m *mongo.Database, x int, y int, squareForBuy int) (bool, error) {
	if squareForBuy <= 0 {
		return false, errors.New("square < 0")
	}

	occupiedLand, err := GetCellOccupiedLand(m, x, y)
	if err != nil {
		return false, errors.New("can't get cell occupied land")
	}

	cell, err := GetCell(m, x, y)
	if err != nil {
		return false, errors.New("can't get cell")
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
		return 0, errors.New("can't get cell owners")
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
		log.Println("Can't get all cells: " + err.Error())
		return cells, err
	}

	err = cursor.All(ctx, &cells)
	return cells, err
}

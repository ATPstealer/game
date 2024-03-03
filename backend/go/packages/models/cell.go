package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type CellMongo struct {
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

func CheckEnoughLandMongo(m *mongo.Database, x int, y int, squareForBuy int) (bool, error) {
	if squareForBuy <= 0 {
		return false, errors.New("square < 0")
	}

	occupiedLand, err := GetCellOccupiedLandMongo(m, x, y)
	if err != nil {
		return false, errors.New("can't get cell occupied land")
	}

	cell, err := GetCellMongo(m, x, y)
	if err != nil {
		return false, errors.New("can't get cell")
	}

	return cell.Square-occupiedLand >= squareForBuy, nil
}

func GetCellMongo(m *mongo.Database, x int, y int) (CellMongo, error) {
	var cell CellMongo
	err := m.Collection("cells").FindOne(context.TODO(),
		bson.M{"x": x, "y": y}).Decode(&cell)
	return cell, err
}

func GetCellOccupiedLandMongo(m *mongo.Database, x int, y int) (int, error) {
	landLords, err := GetCellOwnersMongo(m, x, y)
	if err != nil {
		return 0, errors.New("can't get cell owners")
	}

	occupiedLand := 0
	for _, landLord := range landLords {
		occupiedLand += landLord.Square
	}
	return occupiedLand, nil
}

func AddCivilSavingsMongo(m *mongo.Database, x int, y int, money float64) error {
	_, err := m.Collection("cells").UpdateOne(context.TODO(),
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

func GetAllCellsMongo(m *mongo.Database) ([]CellMongo, error) {
	var cells []CellMongo

	cursor, err := m.Collection("cells").Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Can't get all cells: " + err.Error())
		return cells, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &cells)
	return cells, err
}

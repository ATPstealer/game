package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
)

type Cell struct {
	gorm.Model
	CellName         string  `json:"cellName"`
	X                int     `json:"x"`
	Y                int     `json:"y"`
	SurfaceImagePath string  `json:"surfaceImagePath"`
	Square           int     `json:"square"`
	Pollution        float64 `json:"pollution"`
	Population       float64 `json:"population"`
	CivilSavings     float64 `json:"civilSavings"`
	SpendRate        float64 `json:"SpendRate"`
	Education        float64 `json:"education"`
	Crime            float64 `json:"crime"`
	Medicine         float64 `json:"medicine"`
	AverageSalary    float64 `json:"averageSalary"`
}

func GetCells(db *gorm.DB) []Cell {
	var cells []Cell
	db.Model(&Cell{}).Find(&cells)
	return cells
}

type CellResult struct {
	ID               uint    `json:"id"`
	CellName         string  `json:"cellName"`
	X                int     `json:"x"`
	Y                int     `json:"y"`
	SurfaceImagePath string  `json:"surfaceImagePath"`
	Square           int     `json:"square"`
	Pollution        float64 `json:"pollution"`
	Population       float64 `json:"population"`
	CivilSavings     float64 `json:"civilSavings"`
	SpendRate        float64 `json:"SpendRate"`
	Education        float64 `json:"education"`
	Crime            float64 `json:"crime"`
	Medicine         float64 `json:"medicine"`
	AverageSalary    float64 `json:"averageSalary"`
}

func GetAllCells(db *gorm.DB) []CellResult {
	var cellResults []CellResult
	db.Model(&Cell{}).Find(&cellResults)
	return cellResults
}

func CheckEnoughLand(db *gorm.DB, x int, y int, squareForBuy int) bool {
	if squareForBuy <= 0 {
		return false
	}
	var landLords []LandLord
	res := db.Where("x = ? and y = ?", x, y).Find(&landLords)
	if res.Error != nil {
		log.Println("Can't get Cell Owners: " + res.Error.Error())
	}
	var cell Cell
	res = db.Where("x = ? and y = ?", x, y).First(&cell)
	if res.Error != nil {
		log.Println("Can't get Cell: " + res.Error.Error())
	}
	for _, landLord := range landLords {
		cell.Square -= landLord.Square
	}
	return cell.Square >= squareForBuy
}

func AddCivilSavings(db *gorm.DB, x int, y int, money float64) error {
	var cell Cell
	db.Where("x = ? AND y = ?", x, y).First(&cell)
	cell.CivilSavings += money
	res := db.Save(&cell)
	return res.Error
}

// mongo

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

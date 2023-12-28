package models

import (
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
	ElementarySchool float64 `json:"elementarySchool"`
	HigherSchool     float64 `json:"higherSchool"`
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
	ElementarySchool float64 `json:"elementarySchool"`
	HigherSchool     float64 `json:"higherSchool"`
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

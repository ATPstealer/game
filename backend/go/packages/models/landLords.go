package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
)

type LandLord struct {
	gorm.Model
	UserID uint `json:"userId"`
	Square int  `json:"square"`
	X      int  `json:"x"`
	Y      int  `json:"y"`
}

type BuyLandPayload struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Square int `json:"square"`
}

func BuyLand(db *gorm.DB, userID uint, payload BuyLandPayload) (float64, error) {
	cellLandLords, err := GetCellOwners(db, payload.X, payload.Y)
	if err != nil {
		return 0, err
	}
	var landOccupied int
	for _, landLord := range cellLandLords {
		landOccupied += landLord.Square
	}
	price := 10 * (float64(landOccupied)*2 + 1 + float64(payload.Square)) * float64(payload.Square) / 2

	if !CheckEnoughMoney(db, userID, price) {
		return 0, errors.New("not enough money")
	}
	if !CheckEnoughLand(db, payload.X, payload.Y, payload.Square) {
		return 0, errors.New("not enough land")
	}
	if err := CreateLandLord(db, userID, payload.Square, payload.X, payload.Y); err != nil {
		return 0, err
	}
	if err := AddMoney(db, userID, (-1)*price); err != nil {
		return 0, err
	}
	return price, nil
}

func CreateLandLord(db *gorm.DB, userID uint, square int, x int, y int) error {
	newLandLord := LandLord{
		UserID: userID,
		X:      x,
		Y:      y,
	}
	result := db.Model(&LandLord{}).Where("user_id = ? AND x = ? AND y = ?", userID, x, y).
		First(&newLandLord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newLandLord.Square = square
			db.Create(&newLandLord)
		} else {
			return result.Error
		}
	} else {
		db.Model(&LandLord{}).Where("user_id = ? AND x = ? AND y = ?", userID, x, y).
			Update("Square", gorm.Expr("Square + ?", square))
	}
	return nil
}

type CellLandLordsResult struct {
	NickName string `json:"nickName"`
	Square   int    `json:"square"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

func GetCellOwners(db *gorm.DB, x int, y int) ([]CellLandLordsResult, error) {
	var cellLandLords []CellLandLordsResult
	res := db.Model(&LandLord{}).
		Select("x", "y", "square", "nick_name").
		Where("x = ? and y = ?", x, y).
		Joins("left join users on land_lords.user_id = users.id").
		Scan(&cellLandLords)
	if res.Error != nil {
		log.Println("Can't get Cell Owners: " + res.Error.Error())
	}
	return cellLandLords, res.Error
}

func GetMyLands(db *gorm.DB, userID uint) ([]CellLandLordsResult, error) {
	var myLands []CellLandLordsResult
	res := db.Model(&LandLord{}).
		Select("x", "y", "square", "nick_name").
		Where("users.id", userID).
		Joins("left join users on land_lords.user_id = users.id").
		Scan(&myLands)
	if res.Error != nil {
		log.Println("Can't get Cell Owners: " + res.Error.Error())
	}
	return myLands, res.Error
}

func GetAllLandLords(db *gorm.DB) ([]CellLandLordsResult, error) {
	var allLandLords []CellLandLordsResult
	res := db.Model(&LandLord{}).
		Select("x", "y", "square", "nick_name").
		Joins("left join users on land_lords.user_id = users.id").
		Scan(&allLandLords)
	if res.Error != nil {
		log.Println("Can't get Cell Owners: " + res.Error.Error())
	}
	return allLandLords, res.Error
}

func CheckEnoughLandForBuilding(db *gorm.DB, userID uint, square int, x int, y int) (bool, error) {
	var myLandInCell CellLandLordsResult
	res := db.Model(&LandLord{}).
		Select("square").
		Where("user_id = ? AND x = ? AND y = ?", userID, x, y).
		First(&myLandInCell)
	if res.Error != nil {
		log.Println("Can't get Cell Owners: " + res.Error.Error())
		return false, res.Error
	}
	myBuildingsInCell, err := GetMyBuildingsInCell(db, userID, x, y)
	if err != nil {
		log.Println("Can't get my buildings: " + res.Error.Error())
		return false, res.Error
	}
	freeLand := myLandInCell.Square
	for _, building := range myBuildingsInCell {
		freeLand -= building.Square
	}
	return freeLand >= square, nil
}

// mongo

type LandLordMongo struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	Square int                `json:"square"`
	X      int                `json:"x"`
	Y      int                `json:"y"`
}

// BuyLandMongo Нужно сделать логику покуки
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

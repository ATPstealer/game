package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type BuildingStatus string

const (
	ConstructionStatus    BuildingStatus = "Construction"
	ReadyStatus           BuildingStatus = "Ready"
	ProductionStatus      BuildingStatus = "Production"
	ResourcesNeededStatus BuildingStatus = "ResourcesNeeded"
	StorageNeededStatus   BuildingStatus = "StorageNeeded"
)

type Building struct {
	gorm.Model
	TypeID      uint           `json:"typeId"`
	UserID      uint           `json:"userId"`
	X           int            `json:"x"`
	Y           int            `json:"y"`
	Square      int            `json:"square"`
	Level       int            `json:"level"`
	Status      BuildingStatus `json:"status"`
	WorkStarted *time.Time     `json:"workStarted"`
	WorkEnd     *time.Time     `json:"workEnd"`
	HiringNeeds int            `json:"hiringNeeds"`
	Salary      float64        `json:"salary"`
	Workers     int            `json:"workers"`
	OnStrike    bool           `json:"onStrike"`
}

func GetAllBuildings(db *gorm.DB) ([]Building, error) {
	var buildings []Building
	err := db.Model(&Building{}).Find(&buildings).Error
	if err != nil {
		return nil, err
	}
	return buildings, nil
}

type ConstructBuildingPayload struct {
	TypeID uint
	X      int
	Y      int
	Square int
}

func ConstructBuilding(db *gorm.DB, userID uint, payload ConstructBuildingPayload) error {
	res, err := CheckEnoughLandForBuilding(db, userID, payload.Square, payload.X, payload.Y)
	if err != nil {
		return err
	}
	if !res {
		return errors.New("you don't have enough place in this cell")
	}
	buildingType, err := GetBuildingTypeByID(db, payload.TypeID)
	if err != nil {
		return err
	}
	if !CheckEnoughMoney(db, userID, float64(buildingType.Cost)*float64(payload.Square)) {
		return errors.New("not enough money")
	}
	return CreateBuilding(db, userID, payload, buildingType.Cost)
}

func CreateBuilding(db *gorm.DB, userID uint, payload ConstructBuildingPayload, cost float64) error {
	if err := AddMoney(db, userID, (-1)*cost*float64(payload.Square)); err != nil {
		return err
	}
	buildingType, err := GetBuildingTypeByID(db, payload.TypeID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	now := time.Now()
	end := now.Add(time.Duration(float64(buildingType.BuildTime) * float64(payload.Square)))
	building := Building{
		TypeID:      payload.TypeID,
		UserID:      userID,
		X:           payload.X,
		Y:           payload.Y,
		Square:      payload.Square,
		Level:       1,
		Status:      ConstructionStatus,
		WorkStarted: &now,
		WorkEnd:     &end,
	}
	return db.Create(&building).Error
}

func DestroyBuilding(db *gorm.DB, userID uint, ID uint) error {
	building, err := GetBuildingByID(db, ID)
	if userID != building.UserID && building.UserID != 0 {
		return errors.New("for attempting to destroy someone else's building, inevitable punishment awaits you")
	}
	if err != nil {
		log.Println("Can't destroy building: " + err.Error())
		return err
	}
	db.Delete(&building)
	return nil
}

type MyBuildingsResult struct {
	ID               uint      `json:"id"`
	TypeID           uint      `json:"typeId"`
	Title            string    `json:"title"`
	X                int       `json:"x"`
	Y                int       `json:"y"`
	Square           int       `json:"square"`
	Level            int       `json:"level"`
	Status           string    `json:"status"`
	WorkStarted      time.Time `json:"workStarted"`
	WorkEnd          time.Time `json:"workEnd"`
	HiringNeeds      int       `json:"hiringNeeds"`
	Salary           float64   `json:"salary"`
	Workers          int       `json:"workers"`
	BuildingGroup    string    `json:"buildingGroup"`
	BuildingSubGroup string    `json:"buildingSubGroup"`
	MaxWorkers       int       `json:"maxWorkers"`
	OnStrike         bool      `json:"onStrike"`
}

func GetMyBuildings(db *gorm.DB, userID uint, buildingID uint) ([]MyBuildingsResult, error) {
	var myBuildings []MyBuildingsResult

	query := db.Model(&Building{}).Where("user_id", userID)
	if buildingID != 0 {
		query = query.Where("buildings.id = ?", buildingID)
	}
	res := query.Select("buildings.id", "buildings.type_id", "title", "x", "y", "square", "level",
		"status", "hiring_needs", "salary", "on_strike", "buildings.workers", "buildings.work_started", "buildings.work_end",
		"building_types.building_group", "building_types.building_sub_group", "building_types.workers AS max_workers").
		Joins("left join building_types on buildings.type_id = building_types.id").
		Scan(&myBuildings)

	if res.Error != nil {
		log.Println("Can't get buildings: " + res.Error.Error())
	}
	return myBuildings, res.Error
}

func GetMyBuildingsInCell(db *gorm.DB, userID uint, x int, y int) ([]Building, error) {
	var myBuildings []Building
	res := db.Model(&Building{}).Where("user_id = ? AND x = ? AND y = ?", userID, x, y).Scan(&myBuildings)
	if res.Error != nil {
		log.Println("Can't get buildings: " + res.Error.Error())
	}
	return myBuildings, res.Error
}

type FindBuildingParams struct {
	ID             *uint
	UserID         *uint
	NickName       *string
	X              *int
	Y              *int
	BuildingTypeID *uint
	Limit          *int
	OrderField     *string
	Order          *string
	Page           *int
}

type BuildingsResult struct {
	Title    string `json:"title"`
	TypeID   uint   `json:"typeID"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Square   int    `json:"square"`
	Level    int    `json:"level"`
	Status   string `json:"status"`
	NickName string `json:"nickName"`
}

func GetBuildings(db *gorm.DB, findBuildingParams FindBuildingParams) ([]BuildingsResult, error) {
	var buildings []BuildingsResult
	var request []string
	if findBuildingParams.ID != nil {
		request = append(request, "buildings.id = "+fmt.Sprint(*findBuildingParams.ID))
	}
	if findBuildingParams.UserID != nil {
		request = append(request, "user_id = "+fmt.Sprint(*findBuildingParams.UserID))
	}
	if findBuildingParams.BuildingTypeID != nil {
		request = append(request, "type_id = "+fmt.Sprint(*findBuildingParams.BuildingTypeID))
	}
	if findBuildingParams.X != nil {
		request = append(request, "x = "+fmt.Sprint(*findBuildingParams.X))
	}
	if findBuildingParams.Y != nil {
		request = append(request, "y = "+fmt.Sprint(*findBuildingParams.Y))
	}
	whereString := strings.Join(request, " AND ")
	limit := 10
	if findBuildingParams.Limit != nil {
		limit = *findBuildingParams.Limit
	}
	start := 0
	if findBuildingParams.Page != nil {
		start = (*findBuildingParams.Page - 1) * limit
	}
	order := ""
	if findBuildingParams.OrderField != nil {
		order += *findBuildingParams.OrderField
	}
	if findBuildingParams.Order != nil {
		order += " " + *findBuildingParams.Order
	}

	res := db.Model(&Building{}).Where(whereString).
		Select("title", "type_id", "x", "y", "square", "level", "status", "nick_name").
		Joins("left join building_types on buildings.type_id = building_types.id").
		Joins("left join users on buildings.user_id = users.id").
		Limit(limit).Offset(start).Order(order).
		Scan(&buildings)

	if res.Error != nil {
		log.Println("Can't get buildings: " + res.Error.Error())
	}

	return buildings, nil
}

func GetBuildingByID(db *gorm.DB, buildingID uint) (Building, error) {
	var building Building
	res := db.Model(&Building{}).Where("id = ?", buildingID).First(&building)
	if res.Error != nil {
		log.Println("Can't get building by ID: " + res.Error.Error())
	}
	return building, res.Error
}

func GetAllReadyStorages(db *gorm.DB) ([]Building, error) {
	var storages []Building
	res := db.Model(&Building{}).Where("type_id = ? AND status = ? AND on_strike = ?", 1, ReadyStatus, false).Scan(&storages)
	if res.Error != nil {
		log.Println("Can't get storages: " + res.Error.Error())
	}
	return storages, res.Error
}

type HiringPayload struct {
	BuildingID  uint    `json:"buildingId"`
	Salary      float64 `json:"salary"`
	HiringNeeds int     `json:"hiringNeeds"`
}

func SetHiring(db *gorm.DB, userID uint, payload HiringPayload) error {
	building, err := GetBuildingByID(db, payload.BuildingID)
	if err != nil {
		return err
	}
	if userID != building.UserID && building.UserID != 0 {
		return errors.New("this building doesn't belong to you")
	}
	buildingType, err := GetBuildingTypeByID(db, building.TypeID)
	if err != nil {
		return err
	}
	hiringMax := buildingType.Workers * building.Level * building.Square
	if payload.HiringNeeds > hiringMax {
		return errors.New(fmt.Sprintf("hiring needs more that maximum(%d)", hiringMax))
	}
	building.Salary = payload.Salary
	building.HiringNeeds = payload.HiringNeeds
	db.Save(&building)
	return nil
}

func GetBuildingsForHiring(db *gorm.DB) ([]Building, error) {
	var buildings []Building
	err := db.Model(&Building{}).Where("salary != 0 and hiring_needs != 0").Find(&buildings).Error
	if err != nil {
		return nil, err
	}
	return buildings, nil
}

// mongo

type BuildingMongo struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TypeID      uint               `json:"typeId" bson:"typeId"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	X           int                `json:"x" bson:"x"`
	Y           int                `json:"y" bson:"y"`
	Square      int                `json:"square" bson:"square"`
	Level       int                `json:"level" bson:"level"`
	Status      BuildingStatus     `json:"status" bson:"status"`
	WorkStarted *time.Time         `json:"workStarted" bson:"workStarted"`
	WorkEnd     *time.Time         `json:"workEnd" bson:"workEnd"`
	HiringNeeds int                `json:"hiringNeeds" bson:"hiringNeeds"`
	Salary      float64            `json:"salary" bson:"salary"`
	Workers     int                `json:"workers" bson:"workers"`
	OnStrike    bool               `json:"onStrike" bson:"onStrike"`
}

func ConstructBuildingMongo(m *mongo.Database, userID primitive.ObjectID, payload ConstructBuildingPayload) error {
	enoughLand, err := CheckEnoughLandForBuildingMongo(m, userID, payload.Square, payload.X, payload.Y)
	if !enoughLand {
		return errors.New("not enough land")
	}
	if err != nil {
		return err
	}

	buildingType, err := GetBuildingTypeByIDMongo(m, payload.TypeID)
	if err != nil {
		return err
	}
	if !CheckEnoughMoneyMongo(m, userID, buildingType.Cost*float64(payload.Square)) {
		return errors.New("not enough money")
	}
	return CreateBuildingMongo(m, userID, payload, buildingType)
}

func CreateBuildingMongo(m *mongo.Database, userID primitive.ObjectID, payload ConstructBuildingPayload, buildingType BuildingTypeMongo) error {
	if err := AddMoneyMongo(m, userID, (-1)*buildingType.Cost*float64(payload.Square)); err != nil {
		return err
	}

	now := time.Now()
	end := now.Add(time.Duration(float64(buildingType.BuildTime) * float64(payload.Square)))
	building := BuildingMongo{
		TypeID:      payload.TypeID,
		UserID:      userID,
		X:           payload.X,
		Y:           payload.Y,
		Square:      payload.Square,
		Level:       1,
		Status:      ConstructionStatus,
		WorkStarted: &now,
		WorkEnd:     &end,
	}

	_, err := m.Collection("buildings").InsertOne(context.TODO(), building)
	return err
}

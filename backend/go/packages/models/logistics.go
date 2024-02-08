package models

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"math"
	"time"
)

type Logistic struct {
	gorm.Model
	ResourceTypeID uint      `json:"resourceTypeId"`
	UserID         uint      `json:"userId"`
	Amount         float64   `json:"amount"`
	FromX          int       `json:"fromX"`
	FromY          int       `json:"fromY"`
	ToX            int       `json:"toX"`
	ToY            int       `json:"toY"`
	WorkEnd        time.Time `json:"workEnd"`
}

type LogisticPayload struct {
	ResourceTypeID uint    `json:"resourceTypeId"`
	Amount         float64 `json:"amount"`
	FromX          int     `json:"fromX"`
	FromY          int     `json:"fromY"`
	ToX            int     `json:"toX"`
	ToY            int     `json:"toY"`
}

func StartLogisticJob(db *gorm.DB, userID uint, logisticPayload LogisticPayload) error {
	resource, err := GetMyResourceInCell(db, logisticPayload.ResourceTypeID, userID, logisticPayload.FromX, logisticPayload.FromY)
	if err != nil {
		return err
	}
	if resource.Amount < logisticPayload.Amount {
		return errors.New("not enough resources in this cell")
	}

	resourceType, err := GetResourceTypesByID(db, logisticPayload.ResourceTypeID)
	if err != nil {
		return errors.New("can't get resource type")
	}
	if !CheckEnoughStorage(db, userID, logisticPayload.ToX, logisticPayload.ToY, logisticPayload.Amount*resourceType.Volume) {
		return errors.New("there is not enough storage capacity in the destination sector")
	}

	// FORMULA: logistic
	distance := math.Sqrt(math.Pow(float64(logisticPayload.FromX-logisticPayload.ToX), 2) + math.Pow(float64(logisticPayload.FromY-logisticPayload.ToY), 2))
	price := (resource.Weight + resource.Volume) * distance * logisticPayload.Amount / 1000
	if !CheckEnoughMoney(db, userID, price) {
		return errors.New("not enough money")
	} else {
		if err := AddMoney(db, userID, (-1)*price); err != nil {
			return err
		}
	}

	log.Println(price)
	if err := AddResource(db, resource.ResourceTypeID, userID, logisticPayload.FromX, logisticPayload.FromY, (-1)*logisticPayload.Amount); err != nil {
		return err
	}

	logistic := Logistic{
		ResourceTypeID: logisticPayload.ResourceTypeID,
		UserID:         userID,
		Amount:         logisticPayload.Amount,
		FromX:          logisticPayload.FromX,
		FromY:          logisticPayload.FromY,
		ToX:            logisticPayload.ToX,
		ToY:            logisticPayload.ToY,
		WorkEnd:        time.Now().Add(time.Second * time.Duration(distance*600)), // TODO: come up with a duration
	}
	db.Create(&logistic)

	return nil
}

type LogisticResult struct {
	ID             uint      `json:"id"`
	ResourceTypeID uint      `json:"resourceTypeId"`
	Amount         float64   `json:"amount"`
	FromX          int       `json:"fromX"`
	FromY          int       `json:"fromY"`
	ToX            int       `json:"toX"`
	ToY            int       `json:"toY"`
	WorkEnd        time.Time `json:"workEnd"`
	ResourceName   string    `json:"resourceName"`
	Volume         float64   `json:"volume"`
}

func GetMyLogistics(db *gorm.DB, userID uint) ([]LogisticResult, error) {
	var logistics []LogisticResult
	res := db.Model(&Logistic{}).Where("user_id", userID).
		Select("logistics.id", "resource_type_id", "amount", "from_x", "from_y", "to_x", "to_y", "work_end", "resource_types.name AS resource_name", "volume").
		Joins("left join resource_types on logistics.resource_type_id = resource_types.id").
		Scan(&logistics)

	if res.Error != nil {
		log.Println("Can't get logistics: " + res.Error.Error())
	}
	return logistics, res.Error
}

func GetDestinationVolume(db *gorm.DB, userID uint, toX int, toY int) float64 {
	var volume float64
	res := db.Model(&Logistic{}).Where("user_id = ? AND to_x = ? AND to_y = ?", userID, toX, toY).
		Select("COALESCE(SUM(volume * amount), 0) AS total").
		Joins("left join resource_types on logistics.resource_type_id = resource_types.id").
		Scan(&volume)

	if res.Error != nil {
		log.Println("Can't get logistics: " + res.Error.Error())
	}

	return volume
}

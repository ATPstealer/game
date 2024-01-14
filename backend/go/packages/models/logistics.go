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

func StartLogisticJob(db *gorm.DB, resourceTypeID uint, userID uint, amount float64, fromX int, fromY int, toX int, toY int) error {
	resource, err := GetMyResourceInCell(db, resourceTypeID, userID, fromX, fromY)
	if err != nil {
		return err
	}
	if resource.Amount < amount {
		return errors.New("not enough resources in this cell")
	}

	resourceType, err := GetResourceTypesByID(db, resourceTypeID)
	if err != nil {
		return errors.New("can't get resource type")
	}
	if !CheckEnoughStorage(db, userID, toX, toY, amount*resourceType.Volume) {
		return errors.New("there is not enough storage capacity in the destination sector")
	}

	distance := math.Sqrt(math.Pow(float64(fromX-toX), 2) + math.Pow(float64(fromY-toY), 2))
	price := (resource.Weight + resource.Volume) * float64(distance) * amount / 1000
	if !CheckEnoughMoney(db, userID, price) {
		return errors.New("not enough money")
	} else {
		if err := AddMoney(db, userID, (-1)*price); err != nil {
			return err
		}
	}

	log.Println(price)
	if err := AddResource(db, resource.ResourceTypeID, userID, fromX, fromY, (-1)*amount); err != nil {
		return err
	}

	logistic := Logistic{
		ResourceTypeID: resourceTypeID,
		UserID:         userID,
		Amount:         amount,
		FromX:          fromX,
		FromY:          fromY,
		ToX:            toX,
		ToY:            toY,
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

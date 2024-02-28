package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// mongo

type LogisticMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	FromX          int                `json:"fromX" bson:"fromX"`
	FromY          int                `json:"fromY" bson:"fromY"`
	ToX            int                `json:"toX" bson:"toX"`
	ToY            int                `json:"toY" bson:"toY"`
	WorkEnd        time.Time          `json:"workEnd" bson:"workEnd"`
}

func StartLogisticJobMongo(m *mongo.Database, userID primitive.ObjectID, logisticPayload LogisticPayload) error {
	if !CheckEnoughResourcesMongo(m, logisticPayload.ResourceTypeID, userID,
		logisticPayload.FromX, logisticPayload.FromY, logisticPayload.Amount) {
		return errors.New("not enough resources in this cell")
	}

	resourceType, err := GetResourceTypesByIDMongo(m, logisticPayload.ResourceTypeID)
	if err != nil {
		return errors.New("can't get resource type")
	}
	if !CheckEnoughStorageMongo(m, userID, logisticPayload.ToX, logisticPayload.ToY, logisticPayload.Amount*resourceType.Volume) {
		return errors.New("there is not enough storage capacity in the destination sector")
	}

	// FORMULA: logistic
	distance := math.Sqrt(math.Pow(float64(logisticPayload.FromX-logisticPayload.ToX), 2) + math.Pow(float64(logisticPayload.FromY-logisticPayload.ToY), 2))
	price := (resourceType.Weight + resourceType.Volume) * distance * logisticPayload.Amount / 1000
	if !CheckEnoughMoneyMongo(m, userID, price) {
		return errors.New("not enough money")
	} else {
		if err := AddMoneyMongo(m, userID, (-1)*price); err != nil {
			return err
		}
	}

	log.Println(price)
	if err := AddResourceMongo(m, logisticPayload.ResourceTypeID, userID, logisticPayload.FromX, logisticPayload.FromY, (-1)*logisticPayload.Amount); err != nil {
		return err
	}

	logistic := LogisticMongo{
		ResourceTypeID: logisticPayload.ResourceTypeID,
		UserID:         userID,
		Amount:         logisticPayload.Amount,
		FromX:          logisticPayload.FromX,
		FromY:          logisticPayload.FromY,
		ToX:            logisticPayload.ToX,
		ToY:            logisticPayload.ToY,
		WorkEnd:        time.Now().Add(time.Second * time.Duration(distance*600)),
	}
	_, err = m.Collection("logistics").InsertOne(context.TODO(), &logistic)
	return err
}

func GetMyLogisticsMongo(m *mongo.Database, userID primitive.ObjectID) ([]bson.M, error) {
	matchStage := bson.D{{"$match", bson.M{"userId": userID}}}
	lookupResourceType := bson.D{{"$lookup", bson.D{
		{"from", "resourceTypes"},
		{"localField", "resourceTypeId"},
		{"foreignField", "id"},
		{"as", "resourceType"},
	}}}

	unwindResourceType := bson.D{{"$unwind", bson.D{
		{"path", "$resourceType"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	pipeline := mongo.Pipeline{matchStage, lookupResourceType, unwindResourceType}
	cursor, err := m.Collection("logistics").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var logistics []bson.M
	if err = cursor.All(context.TODO(), &logistics); err != nil {
		log.Println(err)
	}
	return logistics, nil
}

type LogisticResultMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	FromX          int                `json:"fromX" bson:"fromX"`
	FromY          int                `json:"fromY" bson:"fromY"`
	ToX            int                `json:"toX" bson:"toX"`
	ToY            int                `json:"toY" bson:"toY"`
	WorkEnd        time.Time          `json:"workEnd" bson:"workEnd"`
	ResourceType   ResourceType       `json:"resourceType" bson:"resourceType"`
}

func GetDestinationVolumeMongo(m *mongo.Database, userID primitive.ObjectID, toX int, toY int) float64 {
	var volume float64

	filter := bson.M{"userId": userID, "toX": toX, "toY": toY}
	matchStage := bson.D{{"$match", filter}}

	lookupResourceType := bson.D{{"$lookup", bson.D{
		{"from", "resourceTypes"},
		{"localField", "resourceTypeId"},
		{"foreignField", "id"},
		{"as", "resourceType"},
	}}}

	unwindResourceType := bson.D{{"$unwind", bson.D{
		{"path", "$resourceType"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	pipeline := mongo.Pipeline{matchStage, lookupResourceType, unwindResourceType}
	cursor, err := m.Collection("logistics").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get logistics: " + err.Error())
		return 0
	}
	defer cursor.Close(context.TODO())

	var logistics []LogisticResultMongo
	if err = cursor.All(context.TODO(), &logistics); err != nil {
		log.Println(err)
		return 0
	}

	for _, logistic := range logistics {
		volume += logistic.Amount * logistic.ResourceType.Volume
	}

	return volume
}

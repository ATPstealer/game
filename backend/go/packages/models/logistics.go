package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math"
	"time"
)

type Logistic struct {
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

type LogisticPayload struct {
	ResourceTypeID uint    `json:"resourceTypeId"`
	Amount         float64 `json:"amount"`
	FromX          int     `json:"fromX"`
	FromY          int     `json:"fromY"`
	ToX            int     `json:"toX"`
	ToY            int     `json:"toY"`
}

func StartLogisticJob(m *mongo.Database, userID primitive.ObjectID, logisticPayload LogisticPayload) error {
	if !CheckEnoughResources(m, logisticPayload.ResourceTypeID, userID,
		logisticPayload.FromX, logisticPayload.FromY, logisticPayload.Amount) {
		return errors.New("not enough resources in this cell")
	}

	resourceType, err := GetResourceTypesByID(m, logisticPayload.ResourceTypeID)
	if err != nil {
		return errors.New("can't get resource type")
	}
	if !CheckEnoughStorage(m, userID, logisticPayload.ToX, logisticPayload.ToY, logisticPayload.Amount*resourceType.Volume) {
		return errors.New("there is not enough storage capacity in the destination sector")
	}

	// FORMULA: logistic
	distance := math.Sqrt(math.Pow(float64(logisticPayload.FromX-logisticPayload.ToX), 2) + math.Pow(float64(logisticPayload.FromY-logisticPayload.ToY), 2))
	price := (resourceType.Weight + resourceType.Volume) * distance * logisticPayload.Amount / 1000
	if !CheckEnough(m, userID, price) {
		return errors.New("not enough money")
	} else {
		if err := AddMoney(m, userID, (-1)*price); err != nil {
			return err
		}
	}

	log.Println(price)
	if err := AddResource(m, logisticPayload.ResourceTypeID, userID, logisticPayload.FromX, logisticPayload.FromY, (-1)*logisticPayload.Amount); err != nil {
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
		WorkEnd:        time.Now().Add(time.Second * time.Duration(distance*600)),
	}
	_, err = m.Collection("logistics").InsertOne(context.TODO(), &logistic)
	return err
}

func GetMyLogistics(m *mongo.Database, userID primitive.ObjectID) ([]bson.M, error) {
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

type LogisticResult struct {
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

func GetDestinationVolume(m *mongo.Database, userID primitive.ObjectID, toX int, toY int) float64 {
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

	var logistics []LogisticResult
	if err = cursor.All(context.TODO(), &logistics); err != nil {
		log.Println(err)
		return 0
	}

	for _, logistic := range logistics {
		volume += logistic.Amount * logistic.ResourceType.Volume
	}

	return volume
}

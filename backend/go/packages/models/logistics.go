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
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	FromX          int                `json:"fromX" bson:"fromX"`
	FromY          int                `json:"fromY" bson:"fromY"`
	ToX            int                `json:"toX" bson:"toX"`
	ToY            int                `json:"toY" bson:"toY"`
	WorkEnd        time.Time          `json:"workEnd" bson:"workEnd"`
}

type LogisticPayload struct {
	ResourceTypeId uint    `json:"resourceTypeId"`
	Amount         float64 `json:"amount"`
	FromX          int     `json:"fromX"`
	FromY          int     `json:"fromY"`
	ToX            int     `json:"toX"`
	ToY            int     `json:"toY"`
}

func StartLogisticJob(m *mongo.Database, userId primitive.ObjectID, logisticPayload LogisticPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	if !CheckEnoughResources(m, logisticPayload.ResourceTypeId, userId,
		logisticPayload.FromX, logisticPayload.FromY, logisticPayload.Amount) {
		return errors.New("not enough resources in this cell")
	}

	resourceType, err := GetResourceTypesByID(m, logisticPayload.ResourceTypeId)
	if err != nil {
		return errors.New("can't get resource type")
	}
	if !CheckEnoughStorage(m, userId, logisticPayload.ToX, logisticPayload.ToY, logisticPayload.Amount*resourceType.Volume) {
		return errors.New("there is not enough storage capacity in the destination sector")
	}

	// FORMULA: logistics
	distance := math.Sqrt(math.Pow(float64(logisticPayload.FromX-logisticPayload.ToX), 2) + math.Pow(float64(logisticPayload.FromY-logisticPayload.ToY), 2))
	price := (resourceType.Weight + resourceType.Volume) * distance * logisticPayload.Amount / 1000
	if !CheckEnoughMoney(m, userId, price) {
		return errors.New("not enough money")
	} else {
		if err := AddMoney(m, userId, (-1)*price); err != nil {
			return err
		}
	}

	if err := AddResource(m, logisticPayload.ResourceTypeId, userId, logisticPayload.FromX, logisticPayload.FromY, (-1)*logisticPayload.Amount); err != nil {
		return err
	}

	logistic := Logistic{
		ResourceTypeId: logisticPayload.ResourceTypeId,
		UserId:         userId,
		Amount:         logisticPayload.Amount,
		FromX:          logisticPayload.FromX,
		FromY:          logisticPayload.FromY,
		ToX:            logisticPayload.ToX,
		ToY:            logisticPayload.ToY,
		WorkEnd:        time.Now().Add(time.Second * time.Duration(distance*600)),
	}
	_, err = m.Collection("logistics").InsertOne(ctx, &logistic)
	return err
}

type LogisticWithData struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	FromX          int                `json:"fromX" bson:"fromX"`
	FromY          int                `json:"fromY" bson:"fromY"`
	ToX            int                `json:"toX" bson:"toX"`
	ToY            int                `json:"toY" bson:"toY"`
	WorkEnd        time.Time          `json:"workEnd" bson:"workEnd"`
	ResourceType   ResourceType       `json:"resourceType" bson:"resourceType"`
}

func GetMyLogistics(m *mongo.Database, userId primitive.ObjectID) ([]LogisticWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	matchStage := bson.D{{"$match", bson.M{"userId": userId}}}
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
	cursor, err := m.Collection("logistics").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
		return nil, err
	}

	var logistics []LogisticWithData
	if err = cursor.All(ctx, &logistics); err != nil {
		log.Println(err)
	}
	return logistics, nil
}

func GetDestinationVolume(m *mongo.Database, userId primitive.ObjectID, toX int, toY int) float64 {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var volume float64

	filter := bson.M{"userId": userId, "toX": toX, "toY": toY}
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
	cursor, err := m.Collection("logistics").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get logistics: " + err.Error())
		return 0
	}

	var logistics []LogisticWithData
	if err = cursor.All(ctx, &logistics); err != nil {
		log.Println(err)
		return 0
	}

	for _, logistic := range logistics {
		volume += logistic.Amount * logistic.ResourceType.Volume
	}

	return volume
}

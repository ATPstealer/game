package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Resource struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
}

type ResourceWithData struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
	ResourceType   ResourceType       `json:"resourceType" bson:"resourceType"`
}

func GetAllResources(m *mongo.Database) ([]ResourceWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
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
	cursor, err := m.Collection("resources").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
		return nil, err
	}

	var resourcesAndTypes []ResourceWithData
	if err = cursor.All(ctx, &resourcesAndTypes); err != nil {
		log.Println(err)
	}
	return resourcesAndTypes, nil
}

func AddResource(m *mongo.Database, resourceTypeId uint, userId primitive.ObjectID, x int, y int, amount float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("resources").UpdateOne(ctx,
		bson.M{
			"userId":         userId,
			"x":              x,
			"y":              y,
			"resourceTypeId": resourceTypeId,
		},
		bson.M{
			"$inc": bson.M{
				"amount": amount,
			},
			"$setOnInsert": bson.M{
				"userId": userId,
				"x":      x,
				"y":      y,
			},
		},
		options.Update().SetUpsert(true))
	return err
}

func GetMyResources(m *mongo.Database, userId primitive.ObjectID, x *int, y *int) ([]bson.M, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	filter = append(filter, bson.E{Key: "userId", Value: userId})
	if x != nil {
		filter = append(filter, bson.E{Key: "x", Value: *x})
	}
	if y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *y})
	}

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
	cursor, err := m.Collection("resources").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
		return nil, err
	}

	var resources []bson.M
	if err = cursor.All(ctx, &resources); err != nil {
		log.Println(err)
	}
	return resources, nil
}

func GetResourceInCell(m *mongo.Database, resourceTypeID uint, userId primitive.ObjectID, x int, y int) (Resource, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var resource Resource
	err := m.Collection("resources").FindOne(ctx, bson.M{
		"userId":         userId,
		"x":              x,
		"y":              y,
		"resourceTypeId": resourceTypeID,
	}).Decode(&resource)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
	}
	return resource, err
}

func CheckEnoughResources(m *mongo.Database, resourceTypeId uint, userId primitive.ObjectID, x int, y int, amount float64) bool {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var resource Resource
	err := m.Collection("resources").FindOne(ctx, bson.M{
		"userId":         userId,
		"x":              x,
		"y":              y,
		"resourceTypeId": resourceTypeId,
	}).Decode(&resource)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
	}

	return resource.Amount >= amount
}

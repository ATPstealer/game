package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type ResourceMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
}

type ResourceWithTypeMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	Amount         float64            `json:"amount" bson:"amount"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
	ResourceType   ResourceTypeMongo  `json:"resourceType" bson:"resourceType"`
}

func GetAllResourcesMongo(m *mongo.Database) ([]ResourceWithTypeMongo, error) {
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
	cursor, err := m.Collection("resources").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var resourcesAndTypes []ResourceWithTypeMongo
	if err = cursor.All(context.TODO(), &resourcesAndTypes); err != nil {
		log.Println(err)
	}
	return resourcesAndTypes, nil
}

func AddResourceMongo(m *mongo.Database, resourceTypeID uint, userID primitive.ObjectID, x int, y int, amount float64) error {
	_, err := m.Collection("resources").UpdateOne(context.TODO(),
		bson.M{
			"userId":         userID,
			"x":              x,
			"y":              y,
			"resourceTypeId": resourceTypeID,
		},
		bson.M{
			"$inc": bson.M{
				"amount": amount,
			},
			"$setOnInsert": bson.M{
				"userId": userID,
				"x":      x,
				"y":      y,
			},
		},
		options.Update().SetUpsert(true))
	return err
}

func GetMyResourcesMongo(m *mongo.Database, userID primitive.ObjectID, x *int, y *int) ([]bson.M, error) {
	filter := bson.D{}
	filter = append(filter, bson.E{Key: "userId", Value: userID})
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

	// Connect the pipeline stages and execute
	pipeline := mongo.Pipeline{matchStage, lookupResourceType, unwindResourceType}
	cursor, err := m.Collection("resources").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var resources []bson.M
	if err = cursor.All(context.TODO(), &resources); err != nil {
		log.Println(err)
	}
	return resources, nil
}

func GetResourceInCellMongo(m *mongo.Database, resourceTypeID uint, userID primitive.ObjectID, x int, y int) (ResourceMongo, error) {
	var resource ResourceMongo
	err := m.Collection("resources").FindOne(context.TODO(), bson.M{
		"userId":         userID,
		"x":              x,
		"y":              y,
		"resourceTypeId": resourceTypeID,
	}).Decode(&resource)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
	}
	return resource, err
}

func CheckEnoughResourcesMongo(m *mongo.Database, resourceTypeID uint, userID primitive.ObjectID, x int, y int, amount float64) bool {
	var resource ResourceMongo

	err := m.Collection("resources").FindOne(context.TODO(), bson.M{
		"userId":         userID,
		"x":              x,
		"y":              y,
		"resourceTypeId": resourceTypeID,
	}).Decode(&resource)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
	}

	return resource.Amount >= amount
}

package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
)

type Resource struct {
	gorm.Model
	ResourceTypeID uint    `json:"resourceTypeId"`
	UserID         uint    `json:"userId"`
	Amount         float64 `json:"amount"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
}

type ResourceResult struct {
	ID             uint    `json:"id"`
	UserID         uint    `json:"userId"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	Amount         float64 `json:"amount"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	Name           string  `json:"name"`
	Volume         float64 `json:"volume"`
	Weight         float64 `json:"weight"`
}

func GetAllResources(db *gorm.DB) ([]ResourceResult, error) {
	var allResources []ResourceResult
	res := db.Model(&Resource{}).
		Select("resources.id", "user_id", "resource_type_id", "amount", "x", "y", "name", "volume", "weight").
		Joins("left join resource_types on resources.resource_type_id = resource_types.id").
		Scan(&allResources)

	if res.Error != nil {
		log.Println("Can't get resources: " + res.Error.Error())
	}
	return allResources, res.Error
}

func GetMyResources(db *gorm.DB, userID uint, x *int, y *int) ([]ResourceResult, error) {
	var resources []ResourceResult
	query := db.Model(&Resource{}).Where("user_id", userID)
	if x != nil {
		query = query.Where("x = ?", *x)
	}
	if y != nil {
		query = query.Where("y = ?", *y)
	}
	res := query.
		Select("resources.id", "resource_type_id", "amount", "x", "y", "name", "volume", "weight").
		Joins("left join resource_types on resources.resource_type_id = resource_types.id").
		Scan(&resources)

	if res.Error != nil {
		log.Println("Can't get resources: " + res.Error.Error())
	}
	return resources, res.Error
}

func GetMyResourceInCell(db *gorm.DB, resourceTypeID uint, userID uint, x int, y int) (ResourceResult, error) {
	var resource ResourceResult
	res := db.Model(&Resource{}).Where("user_id = ? AND resource_type_id = ? AND X = ? AND Y = ?",
		userID, resourceTypeID, x, y).
		Select("resources.id", "resource_type_id", "amount", "x", "y", "name", "volume", "weight").
		Joins("left join resource_types on resources.resource_type_id = resource_types.id").
		First(&resource)

	if res.Error != nil {
		log.Println("Can't get resources: " + res.Error.Error())
	}
	return resource, res.Error
}

func AddResource(db *gorm.DB, resourceTypeID uint, userID uint, x int, y int, amount float64) error {
	newResource := Resource{
		ResourceTypeID: resourceTypeID,
		UserID:         userID,
		X:              x,
		Y:              y,
	}
	result := db.Model(&Resource{}).Where("resource_type_id =? AND user_id = ? AND x = ? AND y = ?", resourceTypeID, userID, x, y).
		First(&newResource)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newResource.Amount = amount
			db.Create(&newResource)
		} else {
			return result.Error
		}
	} else {
		newResource.Amount += amount
		db.Save(&newResource)
	}
	return nil
}

func CheckEnoughResources(db *gorm.DB, resourceTypeID uint, userID uint, x int, y int, amount float64) bool {
	var resource Resource
	result := db.Model(&Resource{}).Where("resource_type_id =? AND user_id = ? AND x = ? AND y = ?", resourceTypeID, userID, x, y).
		First(&resource)
	if result.Error != nil {
		return false
	}
	return resource.Amount >= amount
}

// mongo

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
	ResourceType   ResourceType       `json:"resourceType" bson:"resourceType"`
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

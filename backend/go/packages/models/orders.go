package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Order struct {
	gorm.Model
	UserID         uint    `json:"userId"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	Amount         float64 `json:"amount"`
	PriceForUnit   float64 `json:"priceForUnit"`
	Sell           bool    `json:"sell"` // true - sell; false - buy
}

type CreateOrderPayload struct {
	ResourceTypeID uint
	X              int
	Y              int
	Amount         float64
	PriceForUnit   float64
	Sell           bool
}

func CreateOrder(db *gorm.DB, userID uint, payload CreateOrderPayload) error {
	if payload.Sell {
		if !CheckEnoughResources(db, payload.ResourceTypeID, userID, payload.X, payload.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		if payload.PriceForUnit < 0 {
			if err := AddMoney(db, userID, payload.PriceForUnit*payload.Amount); err != nil {
				return err
			}
		}
		if err := AddResource(db, payload.ResourceTypeID, userID, payload.X, payload.Y, (-1)*payload.Amount); err != nil {
			return err
		}

	} else {
		if payload.PriceForUnit >= 0 {
			if err := AddMoney(db, userID, (-1)*payload.Amount*payload.PriceForUnit); err != nil {
				return err
			}
		}
	}
	order := Order{
		UserID:         userID,
		X:              payload.X,
		Y:              payload.Y,
		ResourceTypeID: payload.ResourceTypeID,
		Amount:         payload.Amount,
		PriceForUnit:   payload.PriceForUnit,
		Sell:           payload.Sell,
	}
	db.Create(&order)
	return nil
}

type OrderResult struct {
	ID             uint    `json:"id"`
	UserID         uint    `json:"userId"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	ResourceName   string  `json:"resourceName"`
	Amount         float64 `json:"amount"`
	PriceForUnit   float64 `json:"priceForUnit"`
	Sell           bool    `json:"sell"` // true - sell; false - buy
}

func GetMyOrders(db *gorm.DB, userID uint) ([]OrderResult, error) {
	var orderResults []OrderResult
	res := db.Model(&Order{}).Where("user_id = ?", userID).
		Select("orders.id", "user_id", "x", "y", "resource_type_id", "resource_types.name AS resource_name",
			"amount", "price_for_unit", "sell").
		Joins("left join resource_types on orders.resource_type_id = resource_types.id").
		Scan(&orderResults)
	if res.Error != nil {
		log.Println("Can't get buildings: " + res.Error.Error())
	}
	return orderResults, res.Error
}

func CloseMyOrder(db *gorm.DB, userID uint, OrderID uint) error {
	var order Order
	res := db.Model(&Order{}).Where("user_id = ? AND id = ?", userID, OrderID).First(&order)
	if res.Error != nil {
		log.Println("Can't get order: " + res.Error.Error())
		return res.Error
	}
	if order.Sell {
		if err := AddResource(db, order.ResourceTypeID, order.UserID, order.X, order.Y, order.Amount); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoney(db, order.UserID, (-1)*order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	} else {
		if order.PriceForUnit > 0 {
			if err := AddMoney(db, order.UserID, order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	}
	db.Delete(&order)
	return nil
}

type ExecuteOrderPayload struct {
	OrderID uint
	Amount  float64
}

func ExecuteOrder(db *gorm.DB, userID uint, payload ExecuteOrderPayload) error {
	var order Order
	res := db.Model(&Order{}).Where("id = ?", payload.OrderID).First(&order)
	if res.Error != nil {
		log.Println("Can't get order: " + res.Error.Error())
		return res.Error
	}

	if order.Amount < payload.Amount {
		return errors.New("requested quantity is greater than available quantity")
	}

	if order.Sell {
		if err := AddMoney(db, userID, (-1)*payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if err := AddResource(db, order.ResourceTypeID, userID, order.X, order.Y, payload.Amount); err != nil {
			return err
		}
		if order.PriceForUnit > 0 {
			if err := AddMoney(db, order.UserID, payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}

	} else {
		if !CheckEnoughResources(db, order.ResourceTypeID, userID, order.X, order.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		// AddMoney checks enough money if price < 0
		if err := AddMoney(db, userID, payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoney(db, order.UserID, (-1)*payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
		if err := AddResource(db, order.ResourceTypeID, userID, order.X, order.Y, (-1)*payload.Amount); err != nil {
			return err
		}
		if err := AddResource(db, order.ResourceTypeID, order.UserID, order.X, order.Y, payload.Amount); err != nil {
			return err
		}

	}
	if order.Amount == payload.Amount {
		db.Delete(&order)
	} else {
		order.Amount -= payload.Amount
		db.Save(&order)
	}
	return nil
}

type FindOrderParams struct {
	ID             *uint
	UserID         *uint
	X              *int
	Y              *int
	ResourceTypeID *uint
	Sell           *bool
	Limit          *int
	OrderField     *string
	Order          *string
	Page           *int
}

// OrdersResult For GetOrders func
type OrdersResult struct {
	ID             uint    `json:"id"`
	UserID         uint    `json:"userId"`
	NickName       string  `json:"nickName"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ResourceTypeID uint    `json:"resourceTypeId"`
	Amount         float64 `json:"amount"`
	PriceForUnit   float64 `json:"priceForUnit"`
	Sell           bool    `json:"sell"` // true - sell; false - buy
	ResourceName   string  `json:"resourceName"`
	Volume         float64 `json:"volume"`
	Weight         float64 `json:"weight"`
}

func GetOrders(db *gorm.DB, findOrderParams FindOrderParams) ([]OrdersResult, error) {
	var orders []OrdersResult
	var request []string
	if findOrderParams.ID != nil {
		request = append(request, "orders.id = "+fmt.Sprint(*findOrderParams.ID))
	}
	if findOrderParams.UserID != nil {
		request = append(request, "user_id = "+fmt.Sprint(*findOrderParams.UserID))
	}
	if findOrderParams.X != nil {
		request = append(request, "x = "+fmt.Sprint(*findOrderParams.X))
	}
	if findOrderParams.Y != nil {
		request = append(request, "y = "+fmt.Sprint(*findOrderParams.Y))
	}
	if findOrderParams.ResourceTypeID != nil {
		request = append(request, "resource_type_id = "+fmt.Sprint(*findOrderParams.ResourceTypeID))
	}
	if findOrderParams.Sell != nil {
		request = append(request, "sell = "+fmt.Sprint(*findOrderParams.Sell))
	}
	whereString := strings.Join(request, " AND ")

	limit := 10
	if findOrderParams.Limit != nil {
		limit = *findOrderParams.Limit
	}
	start := 0
	if findOrderParams.Page != nil {
		start = (*findOrderParams.Page - 1) * limit
	}
	order := ""
	if findOrderParams.OrderField != nil {
		order += *findOrderParams.OrderField
	}
	if findOrderParams.Order != nil {
		order += " " + *findOrderParams.Order
	}

	res := db.Model(&Order{}).Where(whereString).
		Select("orders.id AS ID", "user_id", "nick_name", "x", "y", "amount", "price_for_unit", "sell",
			"resource_types.id AS ResourceTypeID", "resource_types.name AS ResourceName", "volume", "weight").
		Joins("left join resource_types on orders.resource_type_id = resource_types.id").
		Joins("left join users on orders.user_id = users.id").
		Limit(limit).Offset(start).Order(order).
		Scan(&orders)

	if res.Error != nil {
		log.Println("Can't get orders: " + res.Error.Error())
	}

	return orders, nil
}

// mongo

type OrderMongo struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	Amount         float64            `json:"amount" bson:"amount"`
	PriceForUnit   float64            `json:"priceForUnit" bson:"priceForUnit"`
	Sell           bool               `json:"sell" bson:"sell"` // true - sell; false - buy
}

func CreateOrderMongo(m *mongo.Database, userID primitive.ObjectID, payload OrderMongo) error {
	if payload.Sell {
		if !CheckEnoughResourcesMongo(m, payload.ResourceTypeID, userID, payload.X, payload.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		if payload.PriceForUnit < 0 {
			if err := AddMoneyMongo(m, userID, payload.PriceForUnit*payload.Amount); err != nil {
				return err
			}
		}
		if err := AddResourceMongo(m, payload.ResourceTypeID, userID, payload.X, payload.Y, (-1)*payload.Amount); err != nil {
			return err
		}

	} else {
		if payload.PriceForUnit >= 0 {
			if err := AddMoneyMongo(m, userID, (-1)*payload.Amount*payload.PriceForUnit); err != nil {
				return err
			}
		}
	}

	order := OrderMongo{
		UserID:         userID,
		X:              payload.X,
		Y:              payload.Y,
		ResourceTypeID: payload.ResourceTypeID,
		Amount:         payload.Amount,
		PriceForUnit:   payload.PriceForUnit,
		Sell:           payload.Sell,
	}
	_, err := m.Collection("orders").InsertOne(context.TODO(), &order)
	return err
}

type OrderMongoWithData struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	Amount         float64            `json:"amount" bson:"amount"`
	PriceForUnit   float64            `json:"priceForUnit" bson:"priceForUnit"`
	Sell           bool               `json:"sell" bson:"sell"` // true - sell; false - buy
	ResourceType   ResourceTypeMongo  `json:"resourceType" bson:"resourceType"`
	NickName       string             `json:"nickName" bson:"nickName"`
}

func GetMyOrdersMongo(m *mongo.Database, userID primitive.ObjectID) ([]OrderMongoWithData, error) {
	filter := bson.D{}
	if userID != primitive.NilObjectID {
		filter = append(filter, bson.E{Key: "userId", Value: userID})
	}
	matchStage := bson.D{{"$match", filter}}

	lookupResourceTypes := bson.D{{"$lookup", bson.D{
		{"from", "resourceTypes"},
		{"localField", "resourceTypeId"},
		{"foreignField", "id"},
		{"as", "resourceType"},
	}}}

	unwindResourceTypes := bson.D{{"$unwind", bson.D{
		{"path", "$resourceType"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	lookupUser := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "userId"},
		{"foreignField", "_id"},
		{"as", "user"},
	}}}

	unwindUser := bson.D{{"$unwind", bson.D{
		{"path", "$user"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	project := bson.D{{"$project", bson.D{
		{"id", 1},
		{"userId", 1},
		{"x", 1},
		{"y", 1},
		{"resourceTypeId", 1},
		{"amount", 1},
		{"priceForUnit", 1},
		{"sell", 1},
		{"resourceType.id", 1},
		{"resourceType.name", 1},
		{"resourceType.volume", 1},
		{"resourceType.weight", 1},
		{"resourceType.demand", 1},
		{"resourceType.storeGroup", 1},
		{"nickName", "$user.nickName"},
	}}}

	pipeline := mongo.Pipeline{matchStage, lookupResourceTypes, lookupUser, unwindResourceTypes, unwindUser, project}
	cursor, err := m.Collection("orders").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get orders: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var orders []OrderMongoWithData
	if err = cursor.All(context.TODO(), &orders); err != nil {
		log.Println(err)
	}
	return orders, nil
}

type FindOrderParamsMongo struct {
	ID             *primitive.ObjectID
	UserID         *primitive.ObjectID
	X              *int
	Y              *int
	ResourceTypeID *uint
	Sell           *bool
	Limit          *int
	OrderField     *string
	Order          *int
	Page           *int
}

func GetOrdersMongo(m *mongo.Database, findOrderParams FindOrderParamsMongo) ([]OrderMongoWithData, error) {
	// create filter for match stage
	filter := bson.D{}
	if findOrderParams.ID != nil {
		filter = append(filter, bson.E{Key: "orders._id", Value: *findOrderParams.ID})
	}
	if findOrderParams.UserID != nil {
		filter = append(filter, bson.E{Key: "userId", Value: *findOrderParams.UserID})
	}
	if findOrderParams.X != nil {
		filter = append(filter, bson.E{Key: "x", Value: *findOrderParams.X})
	}
	if findOrderParams.Y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *findOrderParams.Y})
	}
	if findOrderParams.ResourceTypeID != nil {
		filter = append(filter, bson.E{Key: "resourceTypeId", Value: *findOrderParams.ResourceTypeID})
	}
	if findOrderParams.Sell != nil {
		filter = append(filter, bson.E{Key: "sell", Value: *findOrderParams.Sell})
	}
	matchStage := bson.D{{"$match", filter}}

	lookupResourceTypes := bson.D{{"$lookup", bson.D{
		{"from", "resourceTypes"},
		{"localField", "resourceTypeId"},
		{"foreignField", "id"},
		{"as", "resourceType"},
	}}}

	unwindResourceTypes := bson.D{{"$unwind", bson.D{
		{"path", "$resourceType"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	lookupUser := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "userId"},
		{"foreignField", "_id"},
		{"as", "user"},
	}}}

	unwindUser := bson.D{{"$unwind", bson.D{
		{"path", "$user"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	project := bson.D{{"$project", bson.D{
		{"id", 1},
		{"userId", 1},
		{"x", 1},
		{"y", 1},
		{"resourceTypeId", 1},
		{"amount", 1},
		{"priceForUnit", 1},
		{"sell", 1},
		{"resourceType.id", 1},
		{"resourceType.name", 1},
		{"resourceType.volume", 1},
		{"resourceType.weight", 1},
		{"resourceType.demand", 1},
		{"resourceType.storeGroup", 1},
		{"nickName", "$user.nickName"},
	}}}

	sort := bson.D{}

	if findOrderParams.OrderField != nil {
		if findOrderParams.Order != nil {
			sort = append(filter, bson.E{Key: *findOrderParams.OrderField, Value: *findOrderParams.Order})
		} else {
			sort = append(filter, bson.E{Key: *findOrderParams.OrderField, Value: 1})
		}
	}

	sortStage := bson.D{}

	if len(sort) != 0 {
		sortStage = bson.D{{"$sort", sort}}
	} else {
		sortStage = bson.D{{"$sort", bson.D{{"_id", -1}}}}
	}

	limit := 20
	if findOrderParams.Limit != nil {
		limit = *findOrderParams.Limit
	}
	limitStage := bson.D{{"$limit", limit}}

	skipStage := bson.D{{"$skip", 0}}
	if findOrderParams.Page != nil {
		skipStage = bson.D{{"$skip", (*findOrderParams.Page - 1) * limit}}
	}

	pipeline := mongo.Pipeline{matchStage, lookupResourceTypes, lookupUser, unwindResourceTypes,
		unwindUser, project, sortStage, skipStage, limitStage}
	cursor, err := m.Collection("orders").Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Println("Can't get orders: " + err.Error())
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var orders []OrderMongoWithData
	if err = cursor.All(context.TODO(), &orders); err != nil {
		log.Println(err)
	}
	return orders, nil
}

type ExecuteOrderPayloadMongo struct {
	OrderID primitive.ObjectID
	Amount  float64
}

func GetOrderByIDMongo(m *mongo.Database, orderID primitive.ObjectID) (OrderMongo, error) {
	var order OrderMongo
	err := m.Collection("orders").FindOne(context.TODO(),
		bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		log.Println("Can't get building by ID: " + err.Error())
	}
	return order, err
}

func ExecuteOrderMongo(m *mongo.Database, userID primitive.ObjectID, payload ExecuteOrderPayloadMongo) error {
	order, err := GetOrderByIDMongo(m, payload.OrderID)
	if err != nil {
		log.Println("Can't get order: " + err.Error())
		return err
	}

	if order.Amount < payload.Amount {
		return errors.New("requested quantity is greater than available quantity")
	}

	if order.Sell {
		if err := AddMoneyMongo(m, userID, (-1)*payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if err := AddResourceMongo(m, order.ResourceTypeID, userID, order.X, order.Y, payload.Amount); err != nil {
			return err
		}
		if order.PriceForUnit > 0 {
			if err := AddMoneyMongo(m, order.UserID, payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}

	} else {
		if !CheckEnoughResourcesMongo(m, order.ResourceTypeID, userID, order.X, order.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		// AddMoney checks enough money if price < 0
		if err := AddMoneyMongo(m, userID, payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoneyMongo(m, order.UserID, (-1)*payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
		if err := AddResourceMongo(m, order.ResourceTypeID, userID, order.X, order.Y, (-1)*payload.Amount); err != nil {
			return err
		}
		if err := AddResourceMongo(m, order.ResourceTypeID, order.UserID, order.X, order.Y, payload.Amount); err != nil {
			return err
		}

	}

	if order.Amount == payload.Amount {
		_, err := m.Collection("orders").DeleteOne(context.TODO(), bson.M{"_id": order.ID})
		if err != nil {
			return err
		}
	} else {
		update := bson.M{
			"$inc": bson.M{
				"amount": -payload.Amount,
			},
		}
		_, err := m.Collection("orders").UpdateOne(context.TODO(), bson.M{"_id": order.ID}, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseMyOrderMongo(m *mongo.Database, userID primitive.ObjectID, orderID primitive.ObjectID) error {
	order, err := GetOrderByIDMongo(m, orderID)
	if err != nil {
		return err
	}
	if order.UserID != userID {
		return errors.New("you are not the owner of this order")
	}

	if order.Sell {
		if err := AddResourceMongo(m, order.ResourceTypeID, order.UserID, order.X, order.Y, order.Amount); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoneyMongo(m, order.UserID, (-1)*order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	} else {
		if order.PriceForUnit > 0 {
			if err := AddMoneyMongo(m, order.UserID, order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	}
	_, err = m.Collection("orders").DeleteOne(context.TODO(), bson.M{"_id": orderID})
	return err
}

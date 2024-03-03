package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Order struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"userId" bson:"userId"`
	X              int                `json:"x" bson:"x"`
	Y              int                `json:"y" bson:"y"`
	ResourceTypeID uint               `json:"resourceTypeId" bson:"resourceTypeId"`
	Amount         float64            `json:"amount" bson:"amount"`
	PriceForUnit   float64            `json:"priceForUnit" bson:"priceForUnit"`
	Sell           bool               `json:"sell" bson:"sell"` // true - sell; false - buy
}

func CreateOrder(m *mongo.Database, userID primitive.ObjectID, payload Order) error {
	if payload.Sell {
		if !CheckEnoughResources(m, payload.ResourceTypeID, userID, payload.X, payload.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		if payload.PriceForUnit < 0 {
			if err := AddMoney(m, userID, payload.PriceForUnit*payload.Amount); err != nil {
				return err
			}
		}
		if err := AddResource(m, payload.ResourceTypeID, userID, payload.X, payload.Y, (-1)*payload.Amount); err != nil {
			return err
		}

	} else {
		if payload.PriceForUnit >= 0 {
			if err := AddMoney(m, userID, (-1)*payload.Amount*payload.PriceForUnit); err != nil {
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
	ResourceType   ResourceType       `json:"resourceType" bson:"resourceType"`
	NickName       string             `json:"nickName" bson:"nickName"`
}

func GetMyOrders(m *mongo.Database, userID primitive.ObjectID) ([]OrderMongoWithData, error) {
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

type FindOrderParams struct {
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

func GetOrders(m *mongo.Database, findOrderParams FindOrderParams) ([]OrderMongoWithData, error) {
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

type ExecuteOrderPayload struct {
	OrderID primitive.ObjectID
	Amount  float64
}

func GetOrderByID(m *mongo.Database, orderID primitive.ObjectID) (Order, error) {
	var order Order
	err := m.Collection("orders").FindOne(context.TODO(),
		bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		log.Println("Can't get building by ID: " + err.Error())
	}
	return order, err
}

func ExecuteOrder(m *mongo.Database, userID primitive.ObjectID, payload ExecuteOrderPayload) error {
	order, err := GetOrderByID(m, payload.OrderID)
	if err != nil {
		log.Println("Can't get order: " + err.Error())
		return err
	}

	if order.Amount < payload.Amount {
		return errors.New("requested quantity is greater than available quantity")
	}

	if order.Sell {
		if err := AddMoney(m, userID, (-1)*payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if err := AddResource(m, order.ResourceTypeID, userID, order.X, order.Y, payload.Amount); err != nil {
			return err
		}
		if order.PriceForUnit > 0 {
			if err := AddMoney(m, order.UserID, payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}

	} else {
		if !CheckEnoughResources(m, order.ResourceTypeID, userID, order.X, order.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		// AddMoney checks enough money if price < 0
		if err := AddMoney(m, userID, payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoney(m, order.UserID, (-1)*payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
		if err := AddResource(m, order.ResourceTypeID, userID, order.X, order.Y, (-1)*payload.Amount); err != nil {
			return err
		}
		if err := AddResource(m, order.ResourceTypeID, order.UserID, order.X, order.Y, payload.Amount); err != nil {
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

func CloseMyOrder(m *mongo.Database, userID primitive.ObjectID, orderID primitive.ObjectID) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	order, err := GetOrderByID(m, orderID)
	if err != nil {
		return err
	}
	if order.UserID != userID {
		return errors.New("you are not the owner of this order")
	}

	if order.Sell {
		if err := AddResource(m, order.ResourceTypeID, order.UserID, order.X, order.Y, order.Amount); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoney(m, order.UserID, (-1)*order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	} else {
		if order.PriceForUnit > 0 {
			if err := AddMoney(m, order.UserID, order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	}
	_, err = m.Collection("orders").DeleteOne(ctx, bson.M{"_id": orderID})
	return err
}

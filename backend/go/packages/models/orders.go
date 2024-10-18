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
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
	X              int                `json:"x" bson:"x" validate:"required"`
	Y              int                `json:"y" bson:"y" validate:"required"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId" validate:"required"`
	Amount         float64            `json:"amount" bson:"amount" validate:"required"`
	PriceForUnit   float64            `json:"priceForUnit" bson:"priceForUnit" validate:"required"`
	Sell           bool               `json:"sell" bson:"sell" validate:"required"` // true - sell; false - buy
} // @name order

func CreateOrder(m *mongo.Database, userId primitive.ObjectID, payload Order) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	if payload.Sell {
		if !CheckEnoughResources(m, payload.ResourceTypeId, userId, payload.X, payload.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		if payload.PriceForUnit < 0 {
			if err := AddMoney(m, userId, payload.PriceForUnit*payload.Amount); err != nil {
				return err
			}
		}
		if err := AddResource(m, payload.ResourceTypeId, userId, payload.X, payload.Y, (-1)*payload.Amount); err != nil {
			return err
		}

	} else {
		if payload.PriceForUnit >= 0 {
			if err := AddMoney(m, userId, (-1)*payload.Amount*payload.PriceForUnit); err != nil {
				return err
			}
		}
	}

	order := Order{
		UserId:         userId,
		X:              payload.X,
		Y:              payload.Y,
		ResourceTypeId: payload.ResourceTypeId,
		Amount:         payload.Amount,
		PriceForUnit:   payload.PriceForUnit,
		Sell:           payload.Sell,
	}
	_, err := m.Collection("orders").InsertOne(ctx, &order)
	return err
}

type OrderWithData struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
	X              int                `json:"x" bson:"x" validate:"required"`
	Y              int                `json:"y" bson:"y" validate:"required"`
	ResourceTypeId uint               `json:"resourceTypeId" bson:"resourceTypeId" validate:"required"`
	Amount         float64            `json:"amount" bson:"amount" validate:"required"`
	PriceForUnit   float64            `json:"priceForUnit" bson:"priceForUnit" validate:"required"`
	Sell           bool               `json:"sell" bson:"sell" validate:"required"` // true - sell; false - buy
	ResourceType   ResourceType       `json:"resourceType" bson:"resourceType" validate:"required"`
	NickName       string             `json:"nickName" bson:"nickName" validate:"required"`
} // @name orderWithData

func GetMyOrders(m *mongo.Database, userId primitive.ObjectID) ([]OrderWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	if userId != primitive.NilObjectID {
		filter = append(filter, bson.E{Key: "userId", Value: userId})
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
	cursor, err := m.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get orders: " + err.Error())
		return nil, err
	}

	var orders []OrderWithData
	if err = cursor.All(ctx, &orders); err != nil {
		log.Println(err)
	}
	return orders, nil
}

type FindOrderParams struct {
	Id             *primitive.ObjectID
	UserId         *primitive.ObjectID
	X              *int
	Y              *int
	ResourceTypeId *uint
	Sell           *bool
	Limit          *int
	OrderField     *string
	Order          *int
	Page           *int
} // @name findOrderParams

func GetOrders(m *mongo.Database, findOrderParams FindOrderParams) ([]OrderWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	if findOrderParams.Id != nil {
		filter = append(filter, bson.E{Key: "orders._id", Value: *findOrderParams.Id})
	}
	if findOrderParams.UserId != nil {
		filter = append(filter, bson.E{Key: "userId", Value: *findOrderParams.UserId})
	}
	if findOrderParams.X != nil {
		filter = append(filter, bson.E{Key: "x", Value: *findOrderParams.X})
	}
	if findOrderParams.Y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *findOrderParams.Y})
	}
	if findOrderParams.ResourceTypeId != nil {
		filter = append(filter, bson.E{Key: "resourceTypeId", Value: *findOrderParams.ResourceTypeId})
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
	sortStage := bson.D{}

	if findOrderParams.OrderField != nil {
		if findOrderParams.Order != nil {
			sort = append(sort, bson.E{Key: *findOrderParams.OrderField, Value: *findOrderParams.Order})
		} else {
			sort = append(sort, bson.E{Key: *findOrderParams.OrderField, Value: 1})
		}
	}

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
	cursor, err := m.Collection("orders").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get orders: " + err.Error())
		return nil, err
	}

	var orders []OrderWithData
	if err = cursor.All(ctx, &orders); err != nil {
		log.Println(err)
	}
	return orders, nil
}

type ExecuteOrderPayload struct {
	OrderId primitive.ObjectID
	Amount  float64
} // @name executeOrderPayload

func GetOrderByID(m *mongo.Database, orderId primitive.ObjectID) (Order, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var order Order
	err := m.Collection("orders").FindOne(ctx,
		bson.M{"_id": orderId}).Decode(&order)
	if err != nil {
		log.Println("Can't get building by Id: " + err.Error())
	}
	return order, err
}

func ExecuteOrder(m *mongo.Database, userID primitive.ObjectID, payload ExecuteOrderPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	order, err := GetOrderByID(m, payload.OrderId)
	if err != nil {
		log.Println("Can't get order: " + err.Error())
		return err
	}

	if order.Amount < payload.Amount {
		return errors.New("the amount you're requesting exceeds the available quantity")
	}

	if order.Sell {
		if err := AddMoney(m, userID, (-1)*payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if err := AddResource(m, order.ResourceTypeId, userID, order.X, order.Y, payload.Amount); err != nil {
			return err
		}
		if order.PriceForUnit > 0 {
			if err := AddMoney(m, order.UserId, payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}

	} else {
		if !CheckEnoughResources(m, order.ResourceTypeId, userID, order.X, order.Y, payload.Amount) {
			return errors.New("not enough resources in this cell")
		}
		// AddMoney checks enough money if price < 0
		if err := AddMoney(m, userID, payload.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoney(m, order.UserId, (-1)*payload.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
		if err := AddResource(m, order.ResourceTypeId, userID, order.X, order.Y, (-1)*payload.Amount); err != nil {
			return err
		}
		if err := AddResource(m, order.ResourceTypeId, order.UserId, order.X, order.Y, payload.Amount); err != nil {
			return err
		}

	}

	if order.Amount == payload.Amount {
		_, err := m.Collection("orders").DeleteOne(ctx, bson.M{"_id": order.Id})
		if err != nil {
			return err
		}
	} else {
		update := bson.M{
			"$inc": bson.M{
				"amount": -payload.Amount,
			},
		}
		_, err := m.Collection("orders").UpdateOne(ctx, bson.M{"_id": order.Id}, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseMyOrder(m *mongo.Database, userId primitive.ObjectID, orderId primitive.ObjectID) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	order, err := GetOrderByID(m, orderId)
	if err != nil {
		return err
	}
	if order.UserId != userId {
		return errors.New("you are not the owner of this order")
	}

	if order.Sell {
		if err := AddResource(m, order.ResourceTypeId, order.UserId, order.X, order.Y, order.Amount); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoney(m, order.UserId, (-1)*order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	} else {
		if order.PriceForUnit > 0 {
			if err := AddMoney(m, order.UserId, order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
	}
	_, err = m.Collection("orders").DeleteOne(ctx, bson.M{"_id": orderId})
	if err != nil {
		log.Println("Can't close order: " + err.Error())
	}
	return err
}

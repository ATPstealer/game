package models

import (
	"errors"
	"fmt"
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

func ExecuteOrder(db *gorm.DB, userID uint, OrderID uint) error {
	var order Order
	res := db.Model(&Order{}).Where("id = ?", OrderID).First(&order)
	if res.Error != nil {
		log.Println("Can't get order: " + res.Error.Error())
		return res.Error
	}
	if order.Sell {
		if err := AddMoney(db, userID, (-1)*order.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if err := AddResource(db, order.ResourceTypeID, userID, order.X, order.Y, order.Amount); err != nil {
			return err
		}
		if order.PriceForUnit > 0 {
			if err := AddMoney(db, order.UserID, order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}

	} else {
		if !CheckEnoughResources(db, order.ResourceTypeID, userID, order.X, order.Y, order.Amount) {
			return errors.New("not enough resources in this cell")
		}
		// AddMoney checks enough money if price < 0
		if err := AddMoney(db, userID, order.Amount*order.PriceForUnit); err != nil {
			return err
		}
		if order.PriceForUnit < 0 {
			if err := AddMoney(db, order.UserID, (-1)*order.Amount*order.PriceForUnit); err != nil {
				return err
			}
		}
		if err := AddResource(db, order.ResourceTypeID, userID, order.X, order.Y, (-1)*order.Amount); err != nil {
			return err
		}
		if err := AddResource(db, order.ResourceTypeID, order.UserID, order.X, order.Y, order.Amount); err != nil {
			return err
		}

	}
	db.Delete(&order)
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

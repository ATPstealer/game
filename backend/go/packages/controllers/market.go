package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetOrders(c *gin.Context) {
	var findOrdersParams models.FindOrderParams
	if c.Query("id") != "" {
		id, err := include.StrToUInt(c, c.Query("id"))
		if err != nil {
			return
		}
		findOrdersParams.ID = &id
	}
	if c.Query("user_id") != "" {
		userID, err := include.StrToUInt(c, c.Query("user_id"))
		if err != nil {
			return
		}
		findOrdersParams.UserID = &userID
	}
	if c.Query("x") != "" {
		x, err := include.StrToInt(c, c.Query("x"))
		if err != nil {
			return
		}
		findOrdersParams.X = &x
	}
	if c.Query("y") != "" {
		y, err := include.StrToInt(c, c.Query("y"))
		if err != nil {
			return
		}
		findOrdersParams.Y = &y
	}
	if c.Query("resource_type_id") != "" {
		resourceTypeID, err := include.StrToUInt(c, c.Query("resource_type_id"))
		if err != nil {
			return
		}
		findOrdersParams.ResourceTypeID = &resourceTypeID
	}
	if c.Query("sell") != "" {
		sell, err := include.StrToBool(c, c.Query("sell"))
		if err != nil {
			return
		}
		findOrdersParams.Sell = &sell
	}
	if c.Query("limit") != "" {
		limit, err := include.StrToInt(c, c.Query("limit"))
		if err != nil {
			return
		}
		findOrdersParams.Limit = &limit
	}
	if c.Query("order") != "" {
		order := c.Query("order")
		findOrdersParams.Order = &order
	}
	if c.Query("order_field") != "" {
		orderField := c.Query("order_field")
		findOrdersParams.OrderField = &orderField
	}
	if c.Query("page") != "" {
		page, err := include.StrToInt(c, c.Query("page"))
		if err != nil {
			return
		}
		findOrdersParams.Page = &page
	}

	orders, err := models.GetOrders(db.DB, findOrdersParams)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get orders: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": orders})
}

func CreateOrder(c *gin.Context) {
	var createOrderPayload models.CreateOrderPayload
	var err error

	if err = include.GetPayload(c, &createOrderPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	err = models.CreateOrder(db.DB, userID, createOrderPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't create order: " + err.Error()})
		log.Println("Can't create order: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You created order")})
}

func GetMyOrders(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	myOrders, err := models.GetMyOrders(db.DB, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get orders: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myOrders})
}

func CloseMyOrder(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	orderID, err := include.StrToUInt(c, c.Query("order_id"))
	if err != nil {
		return
	}
	err = models.CloseMyOrder(db.DB, userID, orderID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't close order: " + err.Error()})
		log.Println("Can't close order: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You closed order")})
}

func ExecuteOrder(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}

	var payload models.ExecuteOrderPayload
	if err := include.GetPayload(c, &payload); err != nil {
		return
	}

	if payload.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Wrong amount"})
		return
	}

	err = models.ExecuteOrder(db.DB, userID, payload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't execute order: " + err.Error()})
		log.Println("Can't execute order: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You executed order")})
}

// mongo

func CreateOrderMongo(c *gin.Context) {
	var orderPayload models.OrderMongo
	var err error

	if err = include.GetPayload(c, &orderPayload); err != nil {
		return
	}

	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	err = models.CreateOrderMongo(db.M, userID, orderPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't create order: " + err.Error()})
		log.Println("Can't create order: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You created order")})
}

func GetMyOrdersMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}
	myOrders, err := models.GetMyOrdersMongo(db.M, userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get orders: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myOrders})
}

func GetOrdersMongo(c *gin.Context) {
	var findOrdersParams models.FindOrderParamsMongo
	if c.Query("id") != "" {
		id, err := include.StrToPrimObjId(c, c.Query("id"))
		if err != nil {
			return
		}
		findOrdersParams.ID = &id
	}
	if c.Query("user_id") != "" {
		userID, err := include.StrToPrimObjId(c, c.Query("user_id"))
		if err != nil {
			return
		}
		findOrdersParams.UserID = &userID
	}
	if c.Query("x") != "" {
		x, err := include.StrToInt(c, c.Query("x"))
		if err != nil {
			return
		}
		findOrdersParams.X = &x
	}
	if c.Query("y") != "" {
		y, err := include.StrToInt(c, c.Query("y"))
		if err != nil {
			return
		}
		findOrdersParams.Y = &y
	}
	if c.Query("resource_type_id") != "" {
		resourceTypeID, err := include.StrToUInt(c, c.Query("resource_type_id"))
		if err != nil {
			return
		}
		findOrdersParams.ResourceTypeID = &resourceTypeID
	}
	if c.Query("sell") != "" {
		sell, err := include.StrToBool(c, c.Query("sell"))
		if err != nil {
			return
		}
		findOrdersParams.Sell = &sell
	}
	if c.Query("limit") != "" {
		limit, err := include.StrToInt(c, c.Query("limit"))
		if err != nil {
			return
		}
		findOrdersParams.Limit = &limit
	}
	if c.Query("order") != "" {
		order, err := include.StrToInt(c, c.Query("order"))
		if err != nil {
			return
		}
		findOrdersParams.Order = &order
	}
	if c.Query("order_field") != "" {
		orderField := c.Query("order_field")
		findOrdersParams.OrderField = &orderField
	}
	if c.Query("page") != "" {
		page, err := include.StrToInt(c, c.Query("page"))
		if err != nil {
			return
		}
		findOrdersParams.Page = &page
	}

	orders, err := models.GetOrdersMongo(db.M, findOrdersParams)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get orders: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": orders})
}

func ExecuteOrderMongo(c *gin.Context) {
	userID, err := include.GetUserIDFromContextMongo(c)
	if err != nil {
		return
	}

	var payload models.ExecuteOrderPayloadMongo
	if err := include.GetPayload(c, &payload); err != nil {
		return
	}

	if payload.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Wrong amount"})
		return
	}

	err = models.ExecuteOrderMongo(db.M, userID, payload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't execute order: " + err.Error()})
		log.Println("Can't execute order: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You executed order")})
}

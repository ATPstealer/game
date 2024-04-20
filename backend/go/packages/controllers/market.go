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

func CreateOrder(c *gin.Context) {
	var orderPayload models.Order
	var err error

	if err = include.GetPayload(c, &orderPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.CreateOrder(db.M, userId, orderPayload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't create order: " + err.Error()})
		log.Println("Can't create order: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You created order")})
}

func GetMyOrders(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	myOrders, err := models.GetMyOrders(db.M, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get orders: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": myOrders})
}

func GetOrders(c *gin.Context) {
	var findOrdersParams models.FindOrderParams
	if c.Query("id") != "" {
		id, err := include.StrToPrimObjId(c, c.Query("id"))
		if err != nil {
			return
		}
		findOrdersParams.Id = &id
	}
	if c.Query("userId") != "" {
		userID, err := include.StrToPrimObjId(c, c.Query("userId"))
		if err != nil {
			return
		}
		findOrdersParams.UserId = &userID
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
	if c.Query("resourceTypeId") != "" {
		resourceTypeID, err := include.StrToUInt(c, c.Query("resourceTypeId"))
		if err != nil {
			return
		}
		findOrdersParams.ResourceTypeId = &resourceTypeID
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
	if c.Query("orderField") != "" {
		orderField := c.Query("orderField")
		findOrdersParams.OrderField = &orderField
	}
	if c.Query("page") != "" {
		page, err := include.StrToInt(c, c.Query("page"))
		if err != nil {
			return
		}
		findOrdersParams.Page = &page
	}

	orders, err := models.GetOrders(db.M, findOrdersParams)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get orders: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": orders})
}

func CloseMyOrder(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	orderID, err := include.StrToPrimObjId(c, c.Query("order_id"))
	if err != nil {
		return
	}
	err = models.CloseMyOrder(db.M, userId, orderID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't close order: " + err.Error()})
		log.Println("Can't close order: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You closed order")})
}

func ExecuteOrder(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
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

	err = models.ExecuteOrder(db.M, userId, payload)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't execute order: " + err.Error()})
		log.Println("Can't execute order: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": fmt.Sprintf("You executed order")})
}

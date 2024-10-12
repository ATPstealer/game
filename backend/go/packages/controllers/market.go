package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CreateOrder
//
//	@Summary	Create a new market order
//	@Tags		market
//	@Accept		json
//	@Produce	json
//	@Param		orderPayload	body		models.Order	true	"Order payload"
//	@Success	200				{object}	JSONResult
//	@Failure	401				{object}	JSONResult
//	@Failure	500				{object}	JSONResult
//	@Router		/market/order/create [post]
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
		if strings.Contains(err.Error(), "not enough resources in this cell") {
			c.JSON(http.StatusOK, gin.H{"code": 22, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough money") {
			c.JSON(http.StatusOK, gin.H{"code": 24, "text": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -6})
}

// GetMyOrders Get my orders
//
//	@Summary	Get my orders
//	@Tags		market
//	@Produce	json
//	@Success	200	{object}	JSONResult{data=[]models.OrderWithData}
//	@Failure	401	{object}	JSONResult
//	@Failure	500	{object}	JSONResult
//	@Router		/market/order/my [get]
func GetMyOrders(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	myOrders, err := models.GetMyOrders(db.M, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100001, "text": "Can't get orders: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": myOrders})
}

// GetOrders
//
//	@Summary	Fetches orders based on various query parameters
//	@Tags		market
//	@Produce	json
//	@Param		id				query		string	false	"Order ID"
//	@Param		userId			query		string	false	"User ID"
//	@Param		x				query		int		false	"X coordinate"
//	@Param		y				query		int		false	"Y coordinate"
//	@Param		resourceTypeId	query		uint	false	"Resource Type ID"
//	@Param		sell			query		bool	false	"Sell flag"
//	@Param		limit			query		int		false	"Limit number of orders"
//	@Param		order			query		int		false	"Order"
//	@Param		orderField		query		string	false	"Order Field"
//	@Param		page			query		int		false	"Page number"
//	@Success	200				{object}	JSONResult{data=[]models.OrderWithData}
//	@Failure	500				{object}	JSONResult
//	@Router		/orders [get]
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": orders})
}

// CloseMyOrder
//
//	@Summary	Close user's order
//	@Tags		market
//	@Accept		json
//	@Produce	json
//	@Param		orderId	query		string	true	"Order ID"
//	@Success	200		{object}	JSONResult
//	@Failure	401		{object}	JSONResult
//	@Failure	500		{object}	JSONResult
//	@Router		/market/order/close [delete]
func CloseMyOrder(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}
	orderID, err := include.StrToPrimObjId(c, c.Query("orderId"))
	if err != nil {
		return
	}
	err = models.CloseMyOrder(db.M, userId, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": -7})
}

// ExecuteOrder
//
//	@Summary	Partially execute an  order
//	@Tags		market
//	@Accept		json
//	@Produce	json
//	@Param		executeOrderPayload	body		models.ExecuteOrderPayload	true	"Order execution payload"
//	@Success	200					{object}	JSONResult
//	@Failure	401					{object}	JSONResult
//	@Failure	500					{object}	JSONResult
//	@Router		/market/order/execute [post]
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
		c.JSON(http.StatusOK, gin.H{"code": 21})
		return
	}

	err = models.ExecuteOrder(db.M, userId, payload)
	if err != nil {
		if strings.Contains(err.Error(), "the amount you're requesting exceeds the available quantity") {
			c.JSON(http.StatusOK, gin.H{"code": 25, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough money") {
			c.JSON(http.StatusOK, gin.H{"code": 24, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough resources in this cell") {
			c.JSON(http.StatusOK, gin.H{"code": 22, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -8})
}

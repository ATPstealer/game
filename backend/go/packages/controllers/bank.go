package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AddOrDeleteCreditTerm
//
//	@Summary		Add / Change / Delete credit terms in bank contracts
//	@Description	Limit > 0, Rate > 0. For change limit send payload: {"Rate": sameAsExisting, "Rating": sameAsExisting, "Adding": true}
//	@Tags			bank
//	@Accept			json
//	@Produce		json
//	@Param			CreditTermsPayload	body		models.CreditTermsPayload	true	"Credit terms payload"
//	@Success		200					{object}	JSONResult
//	@Failure		401					{object}	JSONResult
//	@Failure		500					{object}	JSONResult
//	@Router			/bank/credit_terms [post]
func AddOrDeleteCreditTerm(c *gin.Context) {
	var creditTermsPayload models.CreditTermsPayload
	if err := include.GetPayload(c, &creditTermsPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.AddOrDeleteCreditTerm(db.M, userId, creditTermsPayload)
	if err != nil {
		if strings.Contains(err.Error(), "this building don't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else if strings.Contains(err.Error(), "parameters must be positive") {
			c.JSON(http.StatusOK, gin.H{"code": 38, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": -16})
}

// GetCreditTerms
//
//	@Summary		Return credit terms
//	@Description	If defined return. Credit term where limit >= in param, rate <= in param, rating <= in param.
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		float64	false	"Credit limit minimum threshold"
//	@Param			rate	query		float64	false	"Credit rate maximum threshold"
//	@Param			rating	query		float64	false	"Credit rating maximum threshold"
//	@Success		200		{object}	JSONResult{data=[]models.CreditTerms}
//	@Failure		500		{object}	JSONResult
//	@Router			/bank/get_credit_terms [get]
func GetCreditTerms(c *gin.Context) {
	var limitPointer, ratePointer, ratingPointer *float64
	if c.Query("limit") != "" {
		limit, err := include.StrToFloat64(c, c.Query("limit"))
		if err != nil {
			return
		}
		limitPointer = &limit
	}
	if c.Query("rate") != "" {
		rate, err := include.StrToFloat64(c, c.Query("rate"))
		if err != nil {
			return
		}
		ratePointer = &rate
	}
	if c.Query("rating") != "" {
		rating, err := include.StrToFloat64(c, c.Query("rating"))
		if err != nil {
			return
		}
		ratingPointer = &rating
	}
	creditTerms, err := models.GetCreditTerms(db.M, limitPointer, ratePointer, ratingPointer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": creditTerms})
}

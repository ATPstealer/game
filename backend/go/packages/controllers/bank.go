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
//	@Description	Limit > 0, Rate > 0
//	@Description	For change limit send payload: {"Rate": sameAsExisting, "Rating": sameAsExisting, "Adding": true}
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

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
//	@Param			creditTermsPayload	body		models.CreditTermsPayload	true	"Credit terms payload"
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
		} else if strings.Contains(err.Error(), "not enough money") {
			c.JSON(http.StatusOK, gin.H{"code": 24, "text": err.Error()})
		} else if strings.Contains(err.Error(), "limit exceeded") {
			c.JSON(http.StatusOK, gin.H{"code": 40, "text": err.Error()})
		} else if strings.Contains(err.Error(), "doesn't have that credit terms") {
			c.JSON(http.StatusOK, gin.H{"code": 41, "text": err.Error()})
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
//	@Tags			bank
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		float64	false	"Credit limit minimum threshold"
//	@Param			rate	query		float64	false	"Credit rate maximum threshold"
//	@Param			rating	query		float64	false	"Credit rating maximum threshold"
//	@Success		200		{object}	JSONResult{data=[]models.CreditTermsWithData}
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

// TakeCredit
//
//	@Summary		Take credit
//	@Description	Get credit in bank. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount":50, "Rate": 0.5, "Rating": -1000000}
//	@Tags			bank
//	@Accept			json
//	@Produce		json
//	@Param			takeCreditPayload	body		models.TakeCreditPayload	true	"Get credit payload"
//	@Success		200					{object}	JSONResult
//	@Failure		401					{object}	JSONResult
//	@Failure		500					{object}	JSONResult
//	@Router			/bank/take_credit [post]
func TakeCredit(c *gin.Context) {
	var takeCreditPayload models.TakeCreditPayload
	if err := include.GetPayload(c, &takeCreditPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.TakeCredit(db.M, userId, takeCreditPayload)

	if err != nil {
		if strings.Contains(err.Error(), "doesn't have that credit terms") {
			c.JSON(http.StatusOK, gin.H{"code": 41, "text": err.Error()})
		} else if strings.Contains(err.Error(), "you don't have enough credit rating") {
			c.JSON(http.StatusOK, gin.H{"code": 42, "text": err.Error()})
		} else if strings.Contains(err.Error(), "amount exceeded") {
			c.JSON(http.StatusOK, gin.H{"code": 43, "text": err.Error()})
		} else if strings.Contains(err.Error(), "you are not a new user") {
			c.JSON(http.StatusOK, gin.H{"code": 44, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": -18})
}

// TakeStateCredit
//
//	@Summary		Take state credit
//	@Description	Get credit from state. Payload example {"buildingId":"670fd64c211de59e1bb8a314", "Amount": 5000}
//	@Tags			bank
//	@Accept			json
//	@Produce		json
//	@Param			takeStateCreditPayload	body		models.TakeStateCreditPayload	true	"Get state credit payload"
//	@Success		200				    	{object}	JSONResult
//	@Failure		401				    	{object}	JSONResult
//	@Failure		500			   			{object}	JSONResult
//	@Router			/bank/take_state_credit [post]
func TakeStateCredit(c *gin.Context) {
	var takeStateCreditPayload models.TakeStateCreditPayload
	if err := include.GetPayload(c, &takeStateCreditPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.TakeStateCredit(db.M, userId, takeStateCreditPayload)

	if err != nil {
		if strings.Contains(err.Error(), "this building don't belong you") {
			c.JSON(http.StatusOK, gin.H{"code": 29, "text": err.Error()})
		} else if strings.Contains(err.Error(), "parameters must be positive") {
			c.JSON(http.StatusOK, gin.H{"code": 38, "text": err.Error()})
		} else if strings.Contains(err.Error(), "limit exceeded") {
			c.JSON(http.StatusOK, gin.H{"code": 40, "text": err.Error()})
		} else if strings.Contains(err.Error(), "doesn't have bank limits") {
			c.JSON(http.StatusOK, gin.H{"code": 40, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": -18})
}

// RepayLoan
//
//	@Summary		Repay loan
//	@Description	Pay off the loan partially or in full. Payload example {"loanId":"674ca2524dfa3a351adbf424", "Amount":122}
//	@Tags			bank
//	@Accept			json
//	@Produce		json
//	@Param			repayLoanPayload       	body		models.RepayLoanPayload	true	"Repay loan payload"
//	@Success		200				    	{object}	JSONResult
//	@Failure		401				    	{object}	JSONResult
//	@Failure		500			   			{object}	JSONResult
//	@Router			/bank/repay_loan [post]
func RepayLoan(c *gin.Context) {
	var repayLoanPayload models.RepayLoanPayload
	if err := include.GetPayload(c, &repayLoanPayload); err != nil {
		return
	}

	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	err = models.RepayLoan(db.M, userId, repayLoanPayload)

	if err != nil {
		if strings.Contains(err.Error(), "you can't pay off someone else's loan") {
			c.JSON(http.StatusOK, gin.H{"code": 46, "text": err.Error()})
		} else if strings.Contains(err.Error(), "not enough money") {
			c.JSON(http.StatusOK, gin.H{"code": 24, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": -18})
}

// GetMyLoans
//
//	@Summary		Get Users Loans
//	@Description	Return all loans connected with userId
//	@Tags			bank
//	@Accept			json
//	@Produce		json
//	@Success		200				    	{object}	JSONResult
//	@Failure		401				    	{object}	JSONResult
//	@Failure		500			   			{object}	JSONResult
//	@Router			/bank/get_loans [get]
func GetMyLoans(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	loans, err := models.GetMyLoans(db.M, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": loans})
}

// DeleteDefaultLoans
//
//	@Summary	Delete Default Loans
//	@Tags		bank
//	@Accept		json
//	@Produce	json
//	@Param		_id	query		string	true	"Loan ID"
//	@Success	200	{object}	JSONResult
//	@Failure	401	{object}	JSONResult
//	@Failure	500	{object}	JSONResult
//	@Router		/bank/delete_loan [delete]
func DeleteDefaultLoans(c *gin.Context) {
	userId, err := include.GetUserIdFromContext(c)
	if err != nil {
		return
	}

	loanId, err := include.StrToPrimObjId(c, c.Query("_id"))
	if err != nil {
		return
	}
	err = models.DeleteDefaultLoans(db.M, userId, loanId)
	if err != nil {
		if strings.Contains(err.Error(), "you can't delete someone else's loan") {
			c.JSON(http.StatusOK, gin.H{"code": 47, "text": err.Error()})
		} else if strings.Contains(err.Error(), "are you crazy? this loan is not default") {
			c.JSON(http.StatusOK, gin.H{"code": 48, "text": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 100001, "text": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": -12})

}

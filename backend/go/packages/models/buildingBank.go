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

type Bank struct {
	LoansAmount         float64 `json:"loansAmount" bson:"loansAmount" validate:"required"`
	LoansLimit          float64 `json:"loansLimit" bson:"loansLimit" validate:"required"`
	LoansAmountNewUsers float64 `json:"loansAmountNewUsers" bson:"loansAmountNewUsers" validate:"required"`
	BorrowedFromState   float64 `json:"borrowedFromState" bson:"borrowedFromState" validate:"required"`
	BorrowedLimit       float64 `json:"borrowedLimit" bson:"borrowedLimit" validate:"required"`
} // @name bank

type CreditTerms struct {
	Limit   float64 `json:"limit" bson:"limit" validate:"required"`
	Rate    float64 `json:"rate" bson:"rate" validate:"required"`
	Rating  float64 `json:"rating" bson:"rating" validate:"required"`
	NewUser bool    `json:"newUser" bson:"newUser" validate:"required"`
} // @name creditTerms

type CreditTermsPayload struct {
	Limit      float64            `json:"limit" validate:"required"`
	Rate       float64            `json:"rate" validate:"required"`
	Rating     float64            `json:"rating" validate:"required"`
	BuildingId primitive.ObjectID `json:"buildingId" validate:"required"`
	Adding     bool               `json:"adding" validate:"required"`
} // @name creditTermsPayload

func AddOrDeleteCreditTerm(m *mongo.Database, userId primitive.ObjectID, payload CreditTermsPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		return err
	}

	if building.UserId != userId {
		return errors.New("this building don't belong you")
	}

	if payload.Limit < 0 || payload.Rate < 0 {
		return errors.New("parameters must be positive")
	}

	if building.Bank == nil {
		return errors.New("bank has no banking limits")
	}

	creditTerm := CreditTerms{
		Limit:   payload.Limit,
		Rate:    payload.Rate,
		Rating:  payload.Rating,
		NewUser: false,
	}

	settings, err := GetSettings(m)
	if err != nil {
		return err
	}

	oldLimit := 0.0
	var amount float64
	var newUserAmount float64
	if payload.Adding {
		if building.CreditTerms == nil {
			building.CreditTerms = &[]CreditTerms{creditTerm}
		} else {
			addCreditTerm(building.CreditTerms, creditTerm, &oldLimit)
		}
		amount = creditTerm.Limit - oldLimit
		newUserAmount = (amount+building.Bank.LoansAmount)*settings["loansForNewUsers"] - building.Bank.LoansAmountNewUsers
		if amount+newUserAmount > 0 && !CheckEnoughMoney(m, userId, amount+newUserAmount) {
			return errors.New("not enough money")
		}
	} else {
		if !removeCreditTerm(building.CreditTerms, creditTerm) {
			return errors.New("doesn't have that credit terms")
		}
		amount = -creditTerm.Limit
		newUserAmount = (amount+building.Bank.LoansAmount)*settings["loansForNewUsers"] - building.Bank.LoansAmountNewUsers
	}

	if amount+building.Bank.LoansAmount > building.Bank.LoansLimit {
		return errors.New("limit exceeded")
	}
	building.Bank.LoansAmount += amount
	building.Bank.LoansAmountNewUsers += newUserAmount

	if err = AddMoney(m, userId, -amount-newUserAmount); err != nil {
		return err
	}

	updateNewUserCreditTermLimit(building.CreditTerms, newUserAmount, settings)

	update := bson.M{
		"$set": bson.M{
			"creditTerms": building.CreditTerms,
			"bank":        building.Bank,
		},
	}
	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": payload.BuildingId}, update)

	return err
}

func addCreditTerm(creditTerms *[]CreditTerms, termToAdd CreditTerms, oldLimit *float64) {
	for i := range *creditTerms {
		if (*creditTerms)[i].Rating == termToAdd.Rating && (*creditTerms)[i].Rate == termToAdd.Rate {
			*oldLimit = (*creditTerms)[i].Limit
			(*creditTerms)[i].Limit = termToAdd.Limit
			return
		}
	}
	*creditTerms = append(*creditTerms, termToAdd)
}

func removeCreditTerm(creditTerms *[]CreditTerms, termToRemove CreditTerms) bool {
	if creditTerms == nil {
		return false
	}

	var found bool
	var updatedCreditTerms []CreditTerms
	for _, ct := range *creditTerms {
		if ct.Limit != termToRemove.Limit || ct.Rate != termToRemove.Rate || ct.Rating != termToRemove.Rating {
			updatedCreditTerms = append(updatedCreditTerms, ct)
		} else {
			found = true
		}
	}
	if !found {
		return false
	}
	*creditTerms = updatedCreditTerms
	return true
}

func updateNewUserCreditTermLimit(creditTerms *[]CreditTerms, amount float64, settings map[string]float64) {
	index := -1
	for i, ct := range *creditTerms {
		if ct.NewUser {
			index = i
		}
	}
	if index != -1 {
		(*creditTerms)[index].Limit += amount
		(*creditTerms)[index].Rate = settings["interestRate"]
		(*creditTerms)[index].Rating = settings["newUserRating"]
	} else {
		*creditTerms = append(*creditTerms, CreditTerms{
			Limit:   amount,
			Rate:    settings["interestRate"],
			Rating:  settings["newUserRating"],
			NewUser: true,
		})
	}
}

type CreditTermsWithData struct {
	BuildingId primitive.ObjectID `json:"buildingId" validate:"required"`
	Limit      float64            `json:"limit" validate:"required"`
	Rate       float64            `json:"rate" validate:"required"`
	Rating     float64            `json:"rating" validate:"required"`
} // @name creditTermsWithData

func GetCreditTerms(m *mongo.Database, limit *float64, rate *float64, rating *float64) ([]CreditTermsWithData, error) {
	var creditTerms []CreditTermsWithData
	banks, err := GetAllReadyBuildingByGroup(m, "Bank")
	if err != nil {
		log.Println(err)
		return creditTerms, err
	}

	for _, bank := range banks {
		if bank.CreditTerms != nil {
			log.Println(bank.CreditTerms)
			for _, ct := range *bank.CreditTerms {
				if limit != nil && !(ct.Limit >= *limit) {
					break
				}
				if rate != nil && !(ct.Rate <= *rate) {
					break
				}
				if rating != nil && !(ct.Rating <= *rating) {
					break
				}
				creditTerms = append(creditTerms, CreditTermsWithData{
					BuildingId: bank.Id,
					Limit:      ct.Limit,
					Rate:       ct.Rate,
					Rating:     ct.Rating,
				})
			}
		}
	}

	return creditTerms, nil
}

func UpdateBankLimits(m *mongo.Database, buildingId primitive.ObjectID, bank Bank) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": buildingId}, bson.M{
		"$set": bson.M{
			"bank": bank,
		},
	})

	if err != nil {
		log.Println("Failed to update building bank: " + err.Error())
	}
}

type TakeCreditPayload struct {
	BuildingId primitive.ObjectID `json:"buildingId" validate:"required"`
	Amount     float64            `json:"amount" validate:"required"`
	Rate       float64            `json:"rate" validate:"required"`
	Rating     float64            `json:"rating" validate:"required"`
} // @name takeCreditPayload

func TakeCredit(m *mongo.Database, userId primitive.ObjectID, payload TakeCreditPayload) error {
	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		return err
	}

	if building.CreditTerms == nil {
		return errors.New("doesn't have that credit terms")
	}

	user, err := GetUserById(m, userId)
	if err != nil {
		return err
	}

	if user.CreditRating < payload.Rating+payload.Amount {
		return errors.New("you don't have enough credit rating")
	}

	index := 0
	found := false
	for i, ct := range *building.CreditTerms {
		if ct.Rating == payload.Rating && ct.Rate == payload.Rate {
			index = i
			found = true
		}
	}
	if !found {
		return errors.New("doesn't have that credit terms")
	}

	if (*building.CreditTerms)[index].NewUser {
		settings, err := GetSettings(m)
		if err != nil {
			return err
		}

		if time.Now().Sub(user.Created) > time.Hour*24*time.Duration(settings["newUserDays"]) {
			return errors.New("you are not a new user")
		}
	}

	if payload.Amount > (*building.CreditTerms)[index].Limit {
		return errors.New("amount exceeded")
	}

	if err = AddMoney(m, userId, payload.Amount); err != nil {
		return err
	}
	if err = IncreaseCreditRating(m, userId, -payload.Amount); err != nil {
		return err
	}

	(*building.CreditTerms)[index].Limit -= payload.Amount
	if (*building.CreditTerms)[index].Limit <= 0 && !(*building.CreditTerms)[index].NewUser {
		*building.CreditTerms = append((*building.CreditTerms)[:index], (*building.CreditTerms)[index+1:]...)
	}

	if err = updateCreditTerms(m, building.Id, building.CreditTerms); err != nil {
		return err
	}

	err = CreateLoan(m, user.Id, building.UserId, payload.Amount, payload.Rate, paying, (*building.CreditTerms)[index].NewUser, false)
	return err
}

func updateCreditTerms(m *mongo.Database, buildingId primitive.ObjectID, creditTerms *[]CreditTerms) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": buildingId}, bson.M{
		"$set": bson.M{
			"creditTerms": creditTerms,
		},
	})

	return err
}

type TakeStateCreditPayload struct {
	BuildingId primitive.ObjectID `json:"buildingId" validate:"required"`
	Amount     float64            `json:"amount" validate:"required"`
} // @name takeStateCreditPayload

func TakeStateCredit(m *mongo.Database, userId primitive.ObjectID, payload TakeStateCreditPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		return err
	}

	if building.UserId != userId {
		return errors.New("this building don't belong you")
	}

	if payload.Amount < 0 {
		return errors.New("parameters must be positive")
	}

	if building.Bank == nil {
		return errors.New("bank has no banking limits")
	}

	if payload.Amount+building.Bank.BorrowedFromState > building.Bank.BorrowedLimit {
		return errors.New("limit exceeded")
	}

	settings, err := GetSettings(m)
	if err != nil {
		return err
	}

	if err = AddMoney(m, userId, payload.Amount); err != nil {
		return err
	}

	if err = CreateLoan(m, userId, primitive.NilObjectID, payload.Amount, settings["interestRate"], paying, false, true); err != nil {
		return err
	}

	building.Bank.BorrowedFromState += payload.Amount
	update := bson.M{
		"$set": bson.M{
			"bank": building.Bank,
		},
	}
	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": payload.BuildingId}, update)

	return err
}

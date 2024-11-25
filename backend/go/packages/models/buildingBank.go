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
	LoansAmount       float64 `json:"loansAmount" bson:"loansAmount" validate:"required"`
	LoansLimit        float64 `json:"loansLimit" bson:"loansLimit" validate:"required"`
	BorrowedFromState float64 `json:"borrowedFromState" bson:"borrowedFromState" validate:"required"`
	BorrowedLimit     float64 `json:"borrowedLimit" bson:"borrowedLimit" validate:"required"`
} // @name bank

type CreditTerms struct {
	Limit  float64 `json:"limit" bson:"limit" validate:"required"`
	Rate   float64 `json:"rate" bson:"rate" validate:"required"`
	Rating float64 `json:"rating" bson:"rating" validate:"required"`
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
		log.Println("Can't get building by Id: " + err.Error())
		return err
	}

	if building.UserId != userId {
		err := errors.New("this building don't belong you")
		return err
	}

	if payload.Limit < 0 || payload.Rate < 0 {
		err := errors.New("parameters must be positive")
		return err
	}

	if building.Bank == nil {
		return errors.New("doesn't have bank limits")
	}

	creditTerms := CreditTerms{
		Limit:  payload.Limit,
		Rate:   payload.Rate,
		Rating: payload.Rating,
	}

	oldLimit := 0.0

	if payload.Adding {
		if building.CreditTerms == nil {
			building.CreditTerms = &[]CreditTerms{creditTerms}
		} else {
			// TODO: make separate func
			found := false
			for i := range *building.CreditTerms {
				if (*building.CreditTerms)[i].Rating == creditTerms.Rating && (*building.CreditTerms)[i].Rate == creditTerms.Rate {
					oldLimit = (*building.CreditTerms)[i].Limit
					(*building.CreditTerms)[i].Limit = creditTerms.Limit
					found = true
					break
				}
			}
			if !found {
				*building.CreditTerms = append(*building.CreditTerms, creditTerms)
			}
		}
	} else {
		updatedCreditTerms := removeCreditTerm(building.CreditTerms, creditTerms)
		building.CreditTerms = &updatedCreditTerms
	}

	var amount float64
	if payload.Adding {
		amount = creditTerms.Limit - oldLimit
		if amount > 0 && !CheckEnoughMoney(m, userId, amount) {
			return errors.New("not enough money")
		}
	} else {
		amount = -creditTerms.Limit
	}

	if amount+building.Bank.LoansAmount > building.Bank.LoansLimit {
		return errors.New("limit exceeded")
	}
	building.Bank.LoansAmount += amount

	err = AddMoney(m, userId, -amount)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"creditTerms": building.CreditTerms,
			"bank":        building.Bank,
		},
	}
	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": payload.BuildingId}, update)

	return err
}

func removeCreditTerm(creditTerms *[]CreditTerms, termToRemove CreditTerms) []CreditTerms {
	if creditTerms == nil {
		return nil
	}

	var updatedCreditTerms []CreditTerms
	for _, ct := range *creditTerms {
		if ct.Limit != termToRemove.Limit || ct.Rate != termToRemove.Rate || ct.Rating != termToRemove.Rating {
			updatedCreditTerms = append(updatedCreditTerms, ct)
		}
	}
	return updatedCreditTerms
}

func GetCreditTerms(m *mongo.Database, limit *float64, rate *float64, rating *float64) ([]CreditTerms, error) {
	var creditTerms []CreditTerms
	banks, err := GetAllReadyBuildingByGroup(m, "Bank")
	log.Println(banks)
	if err != nil {
		log.Println(err)
		return creditTerms, err
	}

	for _, bank := range banks {
		if bank.CreditTerms != nil {
			log.Println(bank.CreditTerms)
			for _, ct := range *bank.CreditTerms {
				log.Println(ct)
				if limit != nil && !(ct.Limit >= *limit) {
					break
				}
				log.Println(ct.Limit)
				if rate != nil && !(ct.Rate <= *rate) {
					break
				}
				if rating != nil && !(ct.Rating <= *rating) {
					break
				}
				creditTerms = append(creditTerms, ct)
			}
		}
	}
	log.Println(creditTerms)
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

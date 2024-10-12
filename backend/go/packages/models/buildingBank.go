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
	LoansAmount       float64 `json:"loansAmount" bson:"loansAmount"`
	LoansLimit        float64 `json:"loansLimit" bson:"loansLimit"`
	BorrowedFromState float64 `json:"borrowedFromState" bson:"borrowedFromState"`
	BorrowedLimit     float64 `json:"borrowedLimit" bson:"borrowedLimit"`
} // @name bank

type CreditTerms struct {
	Limit  float64 `json:"limit" bson:"limit"`
	Rate   float64 `json:"rate" bson:"rate"`
	Rating float64 `json:"rating" bson:"rating"`
} // @name creditTerms

type CreditTermsPayload struct {
	Limit      float64            `json:"limit"`
	Rate       float64            `json:"rate"`
	Rating     float64            `json:"rating"`
	BuildingId primitive.ObjectID `json:"buildingId"`
	Adding     bool               `json:"adding"`
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

	creditTerms := CreditTerms{
		Limit:  payload.Limit,
		Rate:   payload.Rate,
		Rating: payload.Rating,
	}

	if payload.Adding {
		if building.CreditTerms == nil {
			building.CreditTerms = &[]CreditTerms{creditTerms}
		} else {
			// TODO: make separate func
			found := false
			for i := range *building.CreditTerms {
				if (*building.CreditTerms)[i].Rating == creditTerms.Rating && (*building.CreditTerms)[i].Rate == creditTerms.Rate {
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

	update := bson.M{
		"$set": bson.M{
			"creditTerms": building.CreditTerms,
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

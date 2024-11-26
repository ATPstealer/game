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
			addCreditTerm(building.CreditTerms, creditTerms, &oldLimit)
		}
	} else {
		removeCreditTerm(building.CreditTerms, creditTerms)
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

func addCreditTerm(creditTerms *[]CreditTerms, termToAdd CreditTerms, oldLimit *float64) {
	found := false
	for i := range *creditTerms {
		if (*creditTerms)[i].Rating == termToAdd.Rating && (*creditTerms)[i].Rate == termToAdd.Rate {
			*oldLimit = (*creditTerms)[i].Limit
			(*creditTerms)[i].Limit = termToAdd.Limit
			found = true
			break
		}
	}
	if !found {
		*creditTerms = append(*creditTerms, termToAdd)
	}
}

func removeCreditTerm(creditTerms *[]CreditTerms, termToRemove CreditTerms) {
	if creditTerms == nil {
		return
	}

	var updatedCreditTerms []CreditTerms
	for _, ct := range *creditTerms {
		if ct.Limit != termToRemove.Limit || ct.Rate != termToRemove.Rate || ct.Rating != termToRemove.Rating {
			updatedCreditTerms = append(updatedCreditTerms, ct)
		}
	}
	*creditTerms = updatedCreditTerms
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

	if user.CreditRating < payload.Rating {
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

	if payload.Amount > (*building.CreditTerms)[index].Limit {
		return errors.New("amount exceeded")
	}

	err = AddMoney(m, userId, payload.Amount)
	if err != nil {
		return err
	}

	(*building.CreditTerms)[index].Limit -= payload.Amount
	if (*building.CreditTerms)[index].Limit <= 0 {
		*building.CreditTerms = append((*building.CreditTerms)[:index], (*building.CreditTerms)[index+1:]...)
	}

	err = updateCreditTerms(m, building.Id, building.CreditTerms)

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

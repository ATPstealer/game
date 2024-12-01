package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LoanStatus string // @name loanStatus

const (
	paying      LoanStatus = "Paying"
	loanDefault LoanStatus = "LoanDefault"
)

type Loan struct {
	Id             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	BorrowerUserId primitive.ObjectID `json:"borrowerUserId" bson:"borrowerUserId" validate:"required"`
	LenderUserId   primitive.ObjectID `json:"lenderUserId" bson:"lenderUserId" validate:"required"`
	Amount         float64            `json:"amount" bson:"amount" validate:"required"`
	Rate           float64            `json:"interest" bson:"interest" validate:"required"`
	Status         LoanStatus         `json:"status" bson:"status" validate:"required"`
	NewUser        bool               `json:"newUser" bson:"newUser" validate:"required"`
	StateLoan      bool               `json:"stateLoan" bson:"stateLoan" validate:"required"`
} // @name loan

func CreateLoan(m *mongo.Database, borrowerUserId primitive.ObjectID, lenderUserId primitive.ObjectID, amount float64, rate float64, status LoanStatus,
	newUser bool, stateLoan bool) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	loan := Loan{
		BorrowerUserId: borrowerUserId,
		LenderUserId:   lenderUserId,
		Amount:         amount,
		Rate:           rate,
		Status:         status,
		NewUser:        newUser,
		StateLoan:      stateLoan,
	}

	_, err := m.Collection("loans").InsertOne(ctx, loan)
	return err
}

func GetLoanById(m *mongo.Database, id primitive.ObjectID) (Loan, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	var loan Loan
	err := m.Collection("loans").FindOne(ctx, bson.M{"_id": id}).Decode(&loan)
	return loan, err
}

func UpdateLoanAmount(m *mongo.Database, id primitive.ObjectID, amount float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var err error
	if amount == 0 {
		_, err = m.Collection("loans").DeleteOne(ctx, bson.M{"_id": id})
	} else {
		_, err = m.Collection("loans").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"amount": amount}})
	}
	return err
}

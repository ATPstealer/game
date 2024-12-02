package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type LoanStatus string // @name loanStatus

const (
	Paying      LoanStatus = "Paying"
	LoanDefault LoanStatus = "LoanDefault"
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

func createLoan(m *mongo.Database, borrowerUserId primitive.ObjectID, lenderUserId primitive.ObjectID, amount float64, rate float64, status LoanStatus,
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

func getLoanById(m *mongo.Database, id primitive.ObjectID) (Loan, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	var loan Loan
	err := m.Collection("loans").FindOne(ctx, bson.M{"_id": id}).Decode(&loan)
	return loan, err
}

func updateLoanAmount(m *mongo.Database, id primitive.ObjectID, amount float64) error {
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

func GetAllLoans(m *mongo.Database) ([]Loan, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	var loans []Loan
	cursor, err := m.Collection("loans").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &loans)
	return loans, err
}

func UpdateLoanStatus(m *mongo.Database, id primitive.ObjectID, status LoanStatus) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	_, err := m.Collection("loans").UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func GetMyLoans(m *mongo.Database, userId primitive.ObjectID) ([]Loan, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var loans []Loan
	cursor, err := m.Collection("loans").Find(ctx, bson.M{"$or": []bson.M{
		{"borrowerUserId": userId},
		{"lenderUserId": userId},
	}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &loans)
	return loans, err
}

func DeleteDefaultLoans(m *mongo.Database, userId primitive.ObjectID, loanId primitive.ObjectID) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	loan, err := getLoanById(m, loanId)
	if err != nil {
		return err
	}

	if userId != loan.LenderUserId {
		return errors.New("you can't delete someone else's loan")
	}

	if loan.Status != LoanDefault {
		return errors.New("are you crazy? this loan is not default")
	}

	_, err = m.Collection("loans").DeleteOne(ctx, bson.M{"_id": loanId})
	return err
}

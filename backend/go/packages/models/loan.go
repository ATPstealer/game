package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
	Id             primitive.ObjectID `json:"id" bson:"_id" validate:"required"`
	BorrowerUserId primitive.ObjectID `json:"borrowerUserId" bson:"borrowerUserId" validate:"required"`
	LenderUserId   primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
	Amount         float64            `json:"amount" bson:"amount" validate:"required"`
	Rate           float64            `json:"interest" bson:"interest" validate:"required"`
	Status         string             `json:"status" bson:"status" validate:"required"`
	NewUser        bool               `json:"newUser" bson:"newUser" validate:"required"`
}

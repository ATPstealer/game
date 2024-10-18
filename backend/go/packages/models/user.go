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

type User struct {
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	NickName        string             `json:"nickName" bson:"nickName" validate:"required"`
	Email           string             `json:"email" bson:"email" validate:"required"`
	Password        string             `json:"password" bson:"password" validate:"required"`
	Money           float64            `json:"money" bson:"money"`
	Characteristics Characteristics    `json:"characteristics" bson:"characteristics"`
} // @name user

type Characteristics struct {
	Memory       int `json:"memory" bson:"memory" validate:"required"`
	Intelligence int `json:"intelligence" bson:"intelligence" validate:"required"`
	Attention    int `json:"attention" bson:"attention" validate:"required"`
	Wits         int `json:"wits" bson:"wits" validate:"required"`
	Multitasking int `json:"multitasking" bson:"multitasking" validate:"required"`
	Management   int `json:"management" bson:"management" validate:"required"`
	Planning     int `json:"planning" bson:"planning" validate:"required"`
} // @name characteristics

type UserPayload struct {
	NickName string `json:"nickName" validate:"required" validate:"required"`
	Email    string `json:"email" validate:"required" validate:"required"`
	Password string `json:"password" validate:"required" validate:"required"`
	TTL      int    `json:"ttl" validate:"required"`
} // @name userPayload

func CreateUser(m *mongo.Database, nickName string, email string, password string) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	user := User{
		NickName: nickName,
		Email:    email,
		Password: password,
		Money:    1000000,
		Characteristics: Characteristics{
			Memory:       3,
			Intelligence: 3,
			Attention:    3,
			Wits:         3,
			Multitasking: 3,
			Management:   3,
			Planning:     3,
		},
	}
	_, err := m.Collection("users").InsertOne(ctx, user)
	return err
}

func GetUserByNickName(m *mongo.Database, nickName string) (User, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var user User
	err := m.Collection("users").FindOne(ctx, bson.D{{"nickName", nickName}}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserById(m *mongo.Database, userID primitive.ObjectID) (User, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var user User
	err := m.Collection("users").FindOne(ctx, bson.D{{"_id", userID}}).Decode(&user)
	return user, err
}

func GetUserNamesByPrefix(m *mongo.Database, prefix string) ([]string, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	cursor, err := m.Collection("users").Find(ctx, bson.M{"nickName": bson.M{"$regex": prefix}})
	if err != nil {
		return nil, err
	}
	var names []string

	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		names = append(names, user.NickName)
	}

	return names, nil
}

func AddMoney(m *mongo.Database, userId primitive.ObjectID, money float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	user, err := GetUserById(m, userId)
	if err != nil {
		return err
	}
	if user.Money+money < 0 {
		return errors.New("not enough money")
	}
	_, err = m.Collection("users").UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$inc": bson.M{"money": money}})
	return err
}

func CheckEnoughMoney(m *mongo.Database, userID primitive.ObjectID, money float64) bool {
	user, err := GetUserById(m, userID)
	if err != nil {
		log.Println(err)
		return false
	}
	return user.Money >= money
}

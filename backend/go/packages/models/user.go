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
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NickName     string             `json:"nickName" bson:"nickName"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Money        float64            `json:"money" bson:"money"`
	Memory       int                `json:"memory" bson:"memory"`
	Intelligence int                `json:"intelligence" bson:"intelligence"`
	Attention    int                `json:"attention" bson:"attention"`
	Wits         int                `json:"wits" bson:"wits"`
	Multitasking int                `json:"multitasking" bson:"multitasking"`
	Management   int                `json:"management" bson:"management"`
	Planning     int                `json:"planning" bson:"planning"`
}

func CreateUser(m *mongo.Database, nickName string, email string, password string) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	user := User{
		NickName:     nickName,
		Email:        email,
		Password:     password,
		Money:        1000000,
		Memory:       3,
		Intelligence: 3,
		Attention:    3,
		Wits:         3,
		Multitasking: 3,
		Management:   3,
		Planning:     3,
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

func GetUserByID(m *mongo.Database, userID primitive.ObjectID) (User, error) {
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

	for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		names = append(names, user.NickName)
	}

	return names, nil
}

func AddMoney(m *mongo.Database, userID primitive.ObjectID, money float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	user, err := GetUserByID(m, userID)
	if err != nil {
		return err
	}
	if user.Money+money < 0 {
		return errors.New("not enough money")
	}
	_, err = m.Collection("users").UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$inc": bson.M{"money": money}})
	return err
}

func CheckEnough(m *mongo.Database, userID primitive.ObjectID, money float64) bool {
	user, err := GetUserByID(m, userID)
	if err != nil {
		log.Println(err)
		return false
	}
	return user.Money >= money
}

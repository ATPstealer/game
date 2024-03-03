package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UserMongo struct {
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

func CreateUserMongo(m *mongo.Database, nickName string, email string, password string) error {
	userMongo := UserMongo{
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
	_, err := m.Collection("users").InsertOne(context.TODO(), userMongo)
	return err
}

func GetUserByNickNameMongo(m *mongo.Database, nickName string) (UserMongo, error) {
	var user UserMongo
	err := m.Collection("users").FindOne(context.TODO(), bson.D{{"nickName", nickName}}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByIDMongo(m *mongo.Database, userID primitive.ObjectID) (UserMongo, error) {
	var user UserMongo
	err := m.Collection("users").FindOne(context.TODO(), bson.D{{"_id", userID}}).Decode(&user)
	return user, err
}

func GetUserNamesByPrefixMongo(m *mongo.Database, prefix string) ([]string, error) {
	cursor, err := m.Collection("users").Find(context.TODO(), bson.M{"nickName": bson.M{"$regex": prefix}})
	if err != nil {
		return nil, err
	}
	var names []string

	for cursor.Next(context.TODO()) {
		var user UserMongo
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		names = append(names, user.NickName)
	}

	return names, nil
}

func AddMoneyMongo(m *mongo.Database, userID primitive.ObjectID, money float64) error {
	user, err := GetUserByIDMongo(m, userID)
	if err != nil {
		return err
	}
	if user.Money+money < 0 {
		return errors.New("not enough money")
	}
	_, err = m.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": userID}, bson.M{"$inc": bson.M{"money": money}})
	return err
}

func CheckEnoughMoneyMongo(m *mongo.Database, userID primitive.ObjectID, money float64) bool {
	user, err := GetUserByIDMongo(m, userID)
	if err != nil {
		log.Println(err)
		return false
	}
	return user.Money >= money
}

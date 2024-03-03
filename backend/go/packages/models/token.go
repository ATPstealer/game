package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Token struct {
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Token     string             `json:"token" bson:"token"`
	TTL       time.Duration      `json:"ttl" bson:"ttl"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func CreateToken(m *mongo.Database, nickName string) (Token, error) {
	tokenString := make([]byte, 32)
	if _, err := rand.Read(tokenString); err != nil {
		log.Println("Can't create token", err.Error())
	}
	hexToken := hex.EncodeToString(tokenString)

	user, _ := GetUserByNickName(m, nickName)

	token := Token{
		UserID:    user.ID,
		Token:     hexToken,
		TTL:       2592000, // 30d TODO: set it from request
		CreatedAt: time.Now(),
	}
	_, err := m.Collection("tokens").InsertOne(context.TODO(), &token)

	return token, err
}

func GetUserIDByToken(m *mongo.Database, secureToken string) (primitive.ObjectID, error) {
	var token Token
	filter := bson.M{"token": secureToken}
	err := m.Collection("tokens").FindOne(context.TODO(), filter).Decode(&token)

	if err != nil {
		log.Println("Token not found: ", err)
		return primitive.NilObjectID, err
	}
	if time.Now().After(token.CreatedAt.Add(time.Second * token.TTL)) {
		log.Println("token expired")
		return primitive.NilObjectID, errors.New("token expired")
	}
	return token.UserID, nil
}

func Delete(m *mongo.Database, token string) error {
	_, err := m.Collection("tokens").DeleteOne(context.TODO(), bson.M{"token": token})
	return err
}

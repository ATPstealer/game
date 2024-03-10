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
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Token     string             `json:"token" bson:"token"`
	TTL       time.Duration      `json:"ttl" bson:"ttl"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func CreateToken(m *mongo.Database, nickName string) (Token, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	tokenString := make([]byte, 32)
	if _, err := rand.Read(tokenString); err != nil {
		log.Println("Can't create token", err.Error())
	}
	hexToken := hex.EncodeToString(tokenString)

	user, _ := GetUserByNickName(m, nickName)

	token := Token{
		UserId:    user.Id,
		Token:     hexToken,
		TTL:       2592000, // 30d TODO: set it from request
		CreatedAt: time.Now(),
	}
	_, err := m.Collection("tokens").InsertOne(ctx, &token)

	return token, err
}

func GetUserIdByToken(m *mongo.Database, secureToken string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var token Token
	filter := bson.M{"token": secureToken}
	err := m.Collection("tokens").FindOne(ctx, filter).Decode(&token)

	if err != nil {
		log.Println("Token not found: ", err)
		return primitive.NilObjectID, err
	}
	if time.Now().After(token.CreatedAt.Add(time.Second * token.TTL)) {
		log.Println("token expired")
		return primitive.NilObjectID, errors.New("token expired")
	}
	return token.UserId, nil
}

func Delete(m *mongo.Database, token string) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("tokens").DeleteOne(ctx, bson.M{"token": token})
	return err
}

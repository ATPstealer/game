package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"time"
)

type Token struct {
	gorm.Model
	UserID uint          `json:"userId"`
	Token  string        `json:"token"`
	TTL    time.Duration `json:"ttl"`
}

func CreateToken(db *gorm.DB, nickName string) (Token, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		log.Fatalln("Can't create token", err.Error())
	}
	hexToken := hex.EncodeToString(token)

	user, _ := GetUserByNickName(db, nickName)

	log.Println(hexToken)
	newToken := Token{
		UserID: user.ID,
		Token:  hexToken,
		TTL:    2592000, // 30d
	}
	result := db.Create(&newToken)
	return newToken, result.Error
}

func DeleteToken(db *gorm.DB, token string) error {
	return db.Where("token = ?", token).Delete(&Token{}).Error
}

func GetUserIDByToken(db *gorm.DB, secureToken string) (uint, error) {
	var token Token
	result := db.Where("token = ?", secureToken).First(&token)
	if result.Error != nil {
		log.Println("Token not found: ", result.Error)
		return 0, result.Error
	}
	if time.Now().After(token.CreatedAt.Add(time.Second * token.TTL)) {
		log.Println("token expired")
		return 0, errors.New("token expired")
	}
	return token.UserID, nil
}

// mongo

type TokenMongo struct {
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Token     string             `json:"token" bson:"token"`
	TTL       time.Duration      `json:"ttl" bson:"ttl"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func CreateTokenMongo(m *mongo.Database, nickName string) (TokenMongo, error) {
	tokenString := make([]byte, 32)
	if _, err := rand.Read(tokenString); err != nil {
		log.Println("Can't create token", err.Error())
	}
	hexToken := hex.EncodeToString(tokenString)

	user, _ := GetUserByNickNameMongo(m, nickName)

	token := TokenMongo{
		UserID:    user.ID,
		Token:     hexToken,
		TTL:       2592000, // 30d TODO: set it from request
		CreatedAt: time.Now(),
	}
	_, err := m.Collection("tokens").InsertOne(context.TODO(), &token)

	return token, err
}

func GetUserIDByTokenMongo(m *mongo.Database, secureToken string) (primitive.ObjectID, error) {
	var token TokenMongo
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

package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
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

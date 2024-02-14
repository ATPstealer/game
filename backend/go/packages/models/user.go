package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model           // This includes some common fields like ID, CreatedAt, UpdatedAt, and DeletedAt.
	NickName     string  `gorm:"unique" json:"nickName"`
	Email        string  `gorm:"unique" json:"email"`
	Password     string  `json:"password"`
	Money        float64 `json:"money"`
	Memory       int     `json:"memory"`
	Intelligence int     `json:"intelligence"`
	Attention    int     `json:"attention"`
	Wits         int     `json:"wits"`
	Multitasking int     `json:"multitasking"`
	Management   int     `json:"management"`
	Planning     int     `json:"planning"`
}

type UserResult struct {
	ID           uint    `json:"id"`
	NickName     string  `json:"nickName"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	Money        float64 `json:"money"`
	Memory       int     `json:"memory"`
	Intelligence int     `json:"intelligence"`
	Attention    int     `json:"attention"`
	Wits         int     `json:"wits"`
	Multitasking int     `json:"multitasking"`
	Management   int     `json:"management"`
	Planning     int     `json:"planning"`
}

func CreateUser(db *gorm.DB, nickName string, email string, password string) error {
	newUser := User{
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
	res := db.Create(&newUser)
	return res.Error
}

func GetUserData(db *gorm.DB, id uint) (UserResult, error) {
	var user UserResult
	if err := db.Model(&User{}).First(&user, id).Error; err != nil {
		log.Println("User not found: " + err.Error())
		return user, err
	}
	return user, nil
}

func GetUserByID(db *gorm.DB, id uint) (User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		log.Println("User not found: " + err.Error())
		return user, err
	}
	return user, nil
}

func GetUserByNickName(db *gorm.DB, nickName string) (User, error) {
	var user User
	result := db.Where("nick_name = ?", nickName).First(&user)
	if result.Error != nil {
		log.Println("User not found: ", result.Error)
		return user, result.Error
	}
	return user, nil
}

func CheckEnoughMoney(db *gorm.DB, userID uint, money float64) bool {
	user, _ := GetUserByID(db, userID)
	return user.Money >= money
}

func AddMoney(db *gorm.DB, userID uint, money float64) error {
	user, _ := GetUserByID(db, userID)
	if user.Money+money < 0 {
		return errors.New("not enough money")
	}
	user.Money += money
	res := db.Save(&user)
	return res.Error
}

func GetUserNamesByPrefix(db *gorm.DB, prefix string) ([]string, error) {
	var users []User
	res := db.Model(&User{}).Where("nick_name LIKE ?", "%"+prefix+"%").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}
	var names []string
	for _, user := range users {
		names = append(names, user.NickName)
	}
	return names, nil
}

// mongo

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

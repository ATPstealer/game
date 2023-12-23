package models

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model           // This includes some common fields like ID, CreatedAt, UpdatedAt, and DeletedAt.
	NickName     string  `gorm:"unique" json:"nickName"`
	Email        string  `gorm:"unique" json:"email"`
	Password     string  `json:"password"`
	Money        float32 `json:"money"`
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
	Money        float32 `json:"money"`
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

func CheckEnoughMoney(db *gorm.DB, userID uint, money float32) bool {
	user, _ := GetUserByID(db, userID)
	return user.Money >= money
}

func AddMoney(db *gorm.DB, userID uint, money float32) error {
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

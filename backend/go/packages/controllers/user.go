package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateUser(c *gin.Context) {
	user, err := include.GetUserFromRequest(c)
	log.Println(user)
	if err != nil {
		return
	}

	// Creating database row
	err = models.CreateUser(db.DB, user.NickName, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't create user: " + err.Error()})
		log.Println("Can't create user: " + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok"})
}

func Login(c *gin.Context) {
	user, err := include.GetUserFromRequest(c)
	log.Println(user)
	if err != nil {
		return
	}

	// Test user
	var userFondedInDB models.User
	userFondedInDB, err = models.GetUserByNickName(db.DB, user.NickName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't get user: " + err.Error()})
		return
	}

	if userFondedInDB.Password != user.Password {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Wrong password"})
		return
	}

	token, err := models.CreateToken(db.DB, user.NickName)
	if err != nil {
		log.Println("Can't create token: ", err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't create token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": token})
}

func Logout(c *gin.Context) {
	secureToken, err := c.Cookie("secureToken")
	if err != nil {
		return
	}
	if err := models.DeleteToken(db.DB, secureToken); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't delete token: " + err.Error()})
		log.Println("Can't delete token: ", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "You logout"})
}

func GetUserData(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	user, err := models.GetUserData(db.DB, userID)
	if err != nil {
		log.Println("Can't get user ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "text": "Can't get user " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": user})
}

func GetUserNamesByPrefix(c *gin.Context) {
	names, err := models.GetUserNamesByPrefix(db.DB, c.Query("prefix"))
	if err != nil {
		log.Println("Can't get users ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "text": "Can't get users " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "text": "ok", "data": names})
}

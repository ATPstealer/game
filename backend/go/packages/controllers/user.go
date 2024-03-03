package controllers

import (
	"backend/packages/controllers/include"
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func CreateUser(c *gin.Context) {
	user, err := include.GetUserFromRequest(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 5, "text": "Can't get user data from POST body: " + err.Error()})
		return
	}

	err = models.CreateUser(db.M, user.NickName, user.Email, user.Password)

	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error") {
			if strings.Contains(err.Error(), "nickName_1") {
				c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 3, "text": "User with this nickname already exists: " + err.Error(), "values": user})
			} else if strings.Contains(err.Error(), "email_1") {
				c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 4, "text": "User with this email already exists: " + err.Error(), "values": user})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 2, "text": "Can't create user: " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": -1, "values": user})
}

func Login(c *gin.Context) {
	userPayload, err := include.GetUserFromRequest(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 5, "text": "Can't get user data from POST body: " + err.Error()})
		return
	}

	userFondedInDB, err := models.GetUserByNickName(db.M, userPayload.NickName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 6, "text": "Can't get user: " + err.Error()})
		return
	}

	if userFondedInDB.Password != userPayload.Password {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 7, "text": "Wrong password"})
		return
	}

	token, err := models.CreateToken(db.M, userPayload.NickName)
	if err != nil {
		log.Println("Can't create token: ", err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 8, "text": "Can't create token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": token})
}

func GetUserData(c *gin.Context) {
	userID, err := include.GetUserIDFromContext(c)
	if err != nil {
		return
	}
	user, err := models.GetUserByID(db.M, userID)
	if err != nil {
		log.Println("Can't get user ", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "code": 6, "text": "Can't get user " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": user})
}

func Logout(c *gin.Context) {
	secureToken, err := c.Cookie("secureToken")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 9, "test": "Token cookie is required" + err.Error()})
		return
	}
	if err := models.Delete(db.M, secureToken); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 12, "text": "Can't delete token: " + err.Error()})
		log.Println("Can't delete token: ", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "code": -2, "text": "You logout"})
}

func GetUserNamesByPrefix(c *gin.Context) {
	names, err := models.GetUserNamesByPrefix(db.M, c.Query("prefix"))
	if err != nil {
		log.Println("Can't get users ", err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 13, "text": "Can't get users " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "code": 0, "data": names})
}

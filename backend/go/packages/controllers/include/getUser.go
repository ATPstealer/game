package include

import (
	"backend/packages/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
)

func GetUserFromRequest(c *gin.Context) (models.User, error) {
	var user models.User
	var response gin.H

	// Get POST body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response = gin.H{"status": "failed", "text": "Can't get POST request body: " + err.Error()}
		log.Println("Can't get POST request body: " + err.Error())
		c.IndentedJSON(http.StatusNotFound, response)
		return user, err
	}

	// Body JSON to User object. Using same struct
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		response = gin.H{"status": "failed", "text": "Can't parse POST data: " + err.Error()}
		log.Println("Can't parse POST data: " + err.Error())
		c.IndentedJSON(http.StatusNotFound, response)
		return user, err
	}
	log.Println(user.NickName, user.Password)

	return user, nil
}

func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, ok := c.Get("userID")
	if !ok {
		log.Println("UserID not unit")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "text": "UserID not unit"})
		c.Abort()
		return 0, errors.New("userID not unit")
	}
	return userID.(uint), nil
}

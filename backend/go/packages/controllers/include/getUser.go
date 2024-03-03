package include

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"net/http"
)

type UserPayload struct {
	NickName string `json:"nickName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TTL      int    `json:"ttl"`
}

func GetUserFromRequest(c *gin.Context) (UserPayload, error) {
	var user UserPayload

	// Get POST body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Can't get POST request body: " + err.Error())
		return user, err
	}

	// Body JSON to User object. Using same struct
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		log.Println("Can't parse POST data: " + err.Error())
		return user, err
	}

	return user, nil
}

func GetUserIDFromContext(c *gin.Context) (primitive.ObjectID, error) {
	userID, ok := c.Get("userID")
	if !ok {
		log.Println("UserID didn't set")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "text": "UserID didn't set"})
		c.Abort()
		return primitive.NilObjectID, errors.New("userID didn't set")
	}
	return userID.(primitive.ObjectID), nil
}

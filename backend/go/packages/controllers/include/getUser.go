package include

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func GetUserIdFromContext(c *gin.Context) (primitive.ObjectID, error) {
	userId, ok := c.Get("userId")
	if !ok {
		log.Println("Can't get user from session")
		c.JSON(http.StatusUnauthorized, gin.H{"code": 100010, "text": "Can't get user from session"})
		c.Abort()
		return primitive.NilObjectID, errors.New("can't get user from session")
	}
	return userId.(primitive.ObjectID), nil
}

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
		log.Println("UserId didn't set")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "text": "UserId didn't set"})
		c.Abort()
		return primitive.NilObjectID, errors.New("userId didn't set")
	}
	return userId.(primitive.ObjectID), nil
}

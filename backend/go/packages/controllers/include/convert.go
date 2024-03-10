package include

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
)

func StrToInt(c *gin.Context, str string) (int, error) {
	number, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println("Can't parse string as int: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't parse string as int: " + err.Error()})
		c.Abort()
		return 0, err
	}
	return int(number), nil
}

func StrToFloat64(c *gin.Context, str string) (float64, error) {
	number, err := strconv.ParseFloat(str, 32)
	if err != nil {
		log.Println("Can't parse string as float64 : " + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't parse string as float64 : " + err.Error()})
		c.Abort()
		return 0, err
	}
	return float64(number), nil
}

func StrToUInt(c *gin.Context, str string) (uint, error) {
	number, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		log.Println("Can't parse string as uint : " + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't parse string as uint : " + err.Error()})
		c.Abort()
		return 0, err
	}
	return uint(number), nil
}

func StrToBool(c *gin.Context, str string) (bool, error) {
	boolean, err := strconv.ParseBool(str)
	if err != nil {
		log.Println("Can't parse string as boolean : " + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "text": "Can't parse string as boolean : " + err.Error()})
		c.Abort()
		return false, err
	}
	return boolean, nil
}

func StrToPrimObjId(c *gin.Context, str string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		log.Println("Can't parse string as objectID : " + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": "failed", "code": 1, "text": "Can't parse string as uint : " + err.Error()})
		c.Abort()
		return primitive.NilObjectID, err
	}
	return objectID, nil
}

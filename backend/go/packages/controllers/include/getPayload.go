package include

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
)

func GetPayload(c *gin.Context, unmarshalStruct interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Can't get Request body: " + err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 100002, "text": "Can't get Request body: " + err.Error()})
		return err
	}
	err = json.Unmarshal([]byte(body), unmarshalStruct)
	if err != nil {
		log.Printf("Can't unmarshal json: " + err.Error() + " String body: " + string(body))
		c.JSON(http.StatusOK, gin.H{"code": 100003, "text": "Can't unmarshal json: " + err.Error()})
		return err
	}
	return nil
}

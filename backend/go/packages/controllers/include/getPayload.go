package include

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

func GetPayload(c *gin.Context, unmarshalStruct interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100002, "text": "Can't get Request body: " + err.Error()})
		return err
	}
	err = json.Unmarshal([]byte(body), unmarshalStruct)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 100003, "text": "Can't parse POST data: " + err.Error() + string(body)})
		return err
	}
	return nil
}

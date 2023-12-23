package router

import (
	"backend/packages/db"
	"backend/packages/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureToken, err := c.Cookie("secureToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "text": "Token cookie is required"}) // TODO: передалать все так
			c.Abort()
			return
		}
		userID, err := models.GetUserIDByToken(db.DB, secureToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "text": "Token incorrect"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

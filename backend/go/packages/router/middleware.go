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
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "code": 9, "text": "Token cookie is required"})
			c.Abort()
			return
		}
		userId, err := models.GetUserIdByToken(db.M, secureToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "code": 10, "text": "Token incorrect"})
			c.Abort()
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

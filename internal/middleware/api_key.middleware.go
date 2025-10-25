package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	expectedAPIKey := os.Getenv("API_KEY")
	if expectedAPIKey == "" {
		expectedAPIKey = "default"

	}
	return func(c *gin.Context) {
		//get api key
		apiKeyHeader := c.GetHeader("X-API-KEY")
		log.Println("API_KEY Header:", apiKeyHeader)

		if apiKeyHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "API key not provided"})
			return
		}
		if apiKeyHeader != expectedAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		//
		c.Next()

	}

}

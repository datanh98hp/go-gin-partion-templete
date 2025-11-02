<<<<<<< HEAD
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
=======
package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	expectedKey := os.Getenv("API_KEY")
	if expectedKey == "" {
		expectedKey = "secret-key"
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		log.Printf("---X-API-Key : %s", apiKey)
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing X-API-Key"})
			return
		}

		if apiKey != expectedKey {
			log.Println("Invalid API Key")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			return
		}

		ctx.Next()
	}
}
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114

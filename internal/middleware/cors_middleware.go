package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		log.Printf("Origin : %s", origin)

		// Allow the specific origin
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")

		//log.Printf("CORSMiddleware : %v", ctx.Writer.Header())
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

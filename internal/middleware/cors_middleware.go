package middleware

import (
	"net/http"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		fontEndURL := utils.GetEnv("FRONT_END_URL", "http://127.0.0.1:5500")
		// Allow the specific origin
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500,"+fontEndURL)
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

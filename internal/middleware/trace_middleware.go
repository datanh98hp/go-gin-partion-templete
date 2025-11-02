package middleware

import (
	"context"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceId := c.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		// add trace id to original context of Golang
		contextVal := context.WithValue(c.Request.Context(), logger.TraceIdKey, traceId)
		c.Request = c.Request.WithContext(contextVal)

		c.Writer.Header().Set("X-Trace-Id", traceId) // add trace id to header response
		// add trace id to gin context
		c.Set(logger.TraceIdKey, traceId)

		c.Next()
	}
}

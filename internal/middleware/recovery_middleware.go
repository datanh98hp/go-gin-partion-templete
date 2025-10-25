package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(recoveryLogger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if errors := recover(); errors != nil {

				stack := debug.Stack()
				// log error
				recoveryLogger.Error().
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Str("client_ip", c.ClientIP()).
					Str("panic", fmt.Sprintf("%v", errors)).
					Str("stack_at", ExtractFirstAppStackLine(stack)).
					Str("stack", string(stack)).
					Msg("Panic occurred")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Please try again later",
				})
			}
		}()
		// Continue to the next middleware
		c.Next()
	}
}

var stackLineRegex = regexp.MustCompile(`(.+\.go:\d+)`)

// ExtractFirstAppStackLine returns the first line of the stack trace that is relevant to the application
func ExtractFirstAppStackLine(stack []byte) string {
	line := bytes.Split(stack, []byte("\n"))

	for _, l := range line {
		if bytes.Contains(l, []byte(".go")) &&
			!bytes.Contains(l, []byte("/runtime/")) &&
			!bytes.Contains(l, []byte("/debug/")) &&
			!bytes.Contains(l, []byte("/recovery_middleware.go")) {
			cleanLine := strings.TrimSpace(string(l))
			match := stackLineRegex.FindStringSubmatch(cleanLine)
			//fmt.Printf("=========match: %+v", match)
			if len(match) > 0 {
				return match[1]
			}

		}
	}
	return ""
}

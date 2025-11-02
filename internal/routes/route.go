<<<<<<< HEAD
package routes

import (
	"user-management-api/internal/middleware"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	// create logger into file with lumberjack lib
	httpLoger := utils.NewLoggerWithPath("./internal/logs/app.log", "infor")
	recoveryLoger := utils.NewLoggerWithPath("./internal/logs/recovery.log", "error")
	ratelimiterLoger := utils.NewLoggerWithPath("./internal/logs/ratelimit.log", "warning")
	//add middleware
	r.Use(
		middleware.RateLimiterMiddleware(ratelimiterLoger),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLoger),
		middleware.RecoveryMiddleware(recoveryLoger), // Recovery middleware
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
	)
	api_v1 := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api_v1)
	}
}
=======
package routes

import (
	"user-management-api/internal/middleware"
	"user-management-api/internal/utils"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

// RegisterRoutes registers routes into the given gin.Engine.
// It sets up middleware for logging, rate limiting, CORS, tracing, recovery and API key validation.
// It also sets up gzip compression.
// It takes a variable number of Route objects, registers them into the "/api/v1" group and adds them to the given gin.Engine.
func RegisterRoutes(r *gin.Engine, routes ...Route) {
	// create logger into file with lumberjack lib
	httpLoger := utils.NewLoggerWithPath("./internal/logs/app.log", "infor")
	recoveryLoger := utils.NewLoggerWithPath("./internal/logs/recovery.log", "error")
	ratelimiterLoger := utils.NewLoggerWithPath("./internal/logs/ratelimit.log", "warning")
	//add middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(
		middleware.RateLimiterMiddleware(ratelimiterLoger),
		middleware.CORSMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLoger),
		middleware.RecoveryMiddleware(recoveryLoger),
		middleware.ApiKeyMiddleware(),
	)
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	api_v1 := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api_v1)
	}
	// handle url not found
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"error": "Not Found", "path": ctx.Request.URL.Path})
	})
}
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114

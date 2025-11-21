package routes

import (
	"user-management-api/internal/middleware"
	v1 "user-management-api/internal/routes/v1"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

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
func RegisterRoutes(r *gin.Engine, authService auth.TokenService, cacheService cache.RedisCacheService, routes ...Route) {
	// create logger into file with lumberjack lib
	httpLoger := utils.NewLoggerWithPath("http.log", "infor")
	recoveryLoger := utils.NewLoggerWithPath("recovery.log", "error")
	ratelimiterLoger := utils.NewLoggerWithPath("ratelimit.log", "warning")
	//add middleware
	r.Use(
		middleware.RateLimiterMiddleware(ratelimiterLoger),
		middleware.CORSMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLoger),
		middleware.RecoveryMiddleware(recoveryLoger),
		middleware.ApiKeyMiddleware(),
	)
	r.Use(gzip.Gzip(gzip.DefaultCompression)) // gzip data to decrese size data response to fe
	api_v1 := r.Group("/api/v1")

	middleware.InitAuthMiddleware(authService, cacheService)
	protected := api_v1.Group("")
	protected.Use(
		middleware.AuthMiddleware(),
	)

	for _, route := range routes {
		switch route.(type) {
		case *v1.AuthRoute:
			route.Register(api_v1)
		default:
			route.Register(api_v1)
		}
	}
	// handle url not found
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{"error": "Not Found", "path": ctx.Request.URL.Path})
	})
}

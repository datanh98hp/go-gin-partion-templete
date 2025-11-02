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

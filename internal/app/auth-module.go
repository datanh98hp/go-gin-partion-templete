package app

import (
	handlers_v1 "user-management-api/internal/handlers/v1"
	"user-management-api/internal/repositories"
	"user-management-api/internal/routes"
	v1 "user-management-api/internal/routes/v1"
	services_v1 "user-management-api/internal/services/v1"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/mail"
	"user-management-api/pkg/rabbitmq"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx ModulesContext, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService, rabbitMq rabbitmq.RabbitMQService) *AuthModule {
	//initialize the routes
	useRepo := repositories.NewUserRepo(ctx.DB)

	//initialize the services
	authService := services_v1.NewAuthService(useRepo, tokenService, cacheService, mailService, rabbitMq)

	//initialize the handlers
	authHandler := handlers_v1.NewAuthHandler(authService)

	//initialize the routes
	userRoute := v1.NewAuthRoute(authHandler)
	return &AuthModule{
		routes: userRoute,
	}
}

func (um *AuthModule) Route() routes.Route {
	return um.routes
}

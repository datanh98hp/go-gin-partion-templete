package app

import (
	handlers_v1 "user-management-api/internal/handlers/v1"
	"user-management-api/internal/repositories"
	"user-management-api/internal/routes"
	v1 "user-management-api/internal/routes/v1"
	services_v1 "user-management-api/internal/services/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule(ctx ModulesContext) *UserModule {
	//initialize the routes
	useRepo := repositories.NewUserRepo(ctx.DB)

	//initialize the services
	usersService := services_v1.NewUsersService(useRepo, ctx.Redis)

	//initialize the handlers
	usersHandler := handlers_v1.NewUsersHandler(usersService)

	//initialize the routes
	userRoute := v1.NewUserRoute(usersHandler)
	return &UserModule{
		routes: userRoute,
	}
}

func (um *UserModule) Route() routes.Route {
	return um.routes
}

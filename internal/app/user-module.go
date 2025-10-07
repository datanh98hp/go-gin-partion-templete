package app

import (
	"user-management-api/internal/handlers"
	"user-management-api/internal/repositories"
	"user-management-api/internal/routes"
	"user-management-api/internal/services"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	//initialize the routes
	useRepo := repositories.NewUserRepo()

	//initialize the services
	usersService := services.NewUsersService(useRepo)

	//initialize the handlers
	usersHandler := handlers.NewUsersHandler(usersService)

	//initialize the routes
	userRoute := routes.NewUserRoute(usersHandler)
	return &UserModule{
		routes: userRoute,
	}
}

func (um *UserModule) Route() routes.Route {
	return um.routes
}

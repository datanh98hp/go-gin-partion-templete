package app

import (
	"user-management-api/internal/config"
	"user-management-api/internal/routes"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Route() routes.Route
}

type Application struct {
	Config *config.Config
	Router *gin.Engine
}

func NewApplication(cfg *config.Config) *Application {
	r := gin.Default()
	modules := []Module{
		NewUserModule(),
	}
	routes.RegisterRoutes(r, getModuleRoutes(modules)...) // Register the routes

	return &Application{
		Config: cfg,
		Router: r,
	}

}

func (r *Application) Run() error {

	return r.Router.Run(r.Config.ServerAddress)
}

func getModuleRoutes(modules []Module) []routes.Route {
	routes := make([]routes.Route, len(modules))
	for i, module := range modules {
		routes[i] = module.Route()
	}
	return routes
}

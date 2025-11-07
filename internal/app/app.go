package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/routes"
	"user-management-api/internal/validations"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Route() routes.Route
}
type ModulesContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}
type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Modules []Module
}

func NewApplication(cfg *config.Config) *Application {

	// Customize Recovery middleware to handle panic and return JSON response

	r := gin.Default()

	//postgres
	if err := db.InitializeDatabase(); err != nil { //db.InitializeDatabase()
		log.Fatalf("Error initializing database: %v", err)
	}
	//Redis

	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)

	// create modules context
	ctx := ModulesContext{
		DB:    db.DB,
		Redis: redisClient,
	}
	//Call validator
	if err := validations.InitValidator(); err != nil {
		log.Fatalf("Error initializing validator: %v", err)
	}
	modules := []Module{
		/// add modules
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheRedisService),
	}
	routes.RegisterRoutes(r, tokenService, cacheRedisService, getModuleRoutes(modules)...) // Register the routes

	return &Application{
		Config:  cfg,
		Router:  r,
		Modules: modules,
	}

}

func (a *Application) Run() error {

	srv := &http.Server{
		Addr:    "" + a.Config.ServerAddress,
		Handler: a.Router,
	}

	/// add channel to listen for interrupt or terminate signal from OS / kill -9
	quit := make(chan os.Signal, 1)
	// Listen for interrupt signals : termination signal
	// syscall.SIGINT : sent by Ctrl+C
	// syscall.SIGTERM : default signal sent by "kill" command
	// syscall.SIGHUP : terminal closed
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// another goroutine to listen for the signal
	go func() {
		log.Printf("Server is running at %s", a.Config.ServerAddress)
		// Start serv
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %s\n", err)
		}
	}()
	<-quit // wait here until we get the signal
	log.Println("Shutting down server...")
	// block until we receive our signal.
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//shutdown the application
	if err := srv.Shutdown(context); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
	log.Println("Server exited gracefully...")
	return nil
	//return a.Router.Run(":" + a.Config.ServerAddress)
}

func getModuleRoutes(modules []Module) []routes.Route {
	routes := make([]routes.Route, len(modules))
	for i, module := range modules {
		routes[i] = module.Route()
	}
	return routes
}

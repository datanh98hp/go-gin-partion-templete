package main

import (
	"user-management-api/internal/app"
	"user-management-api/internal/config"
)

func main() {
	// Run the application
	//initialize the config
	cfg := config.NewConfig()
	//initialize theapplicationn
	application := app.NewApplication(cfg)

	//start server
	if err := application.Run(); err != nil {
		panic(err)
	}

}

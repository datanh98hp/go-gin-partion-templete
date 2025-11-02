<<<<<<< HEAD
package main

import (
	"log"
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
		// panic(err)
		log.Fatal(err)
	}

}
=======
package main

import (
	"log"
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
		// panic(err)
		log.Fatal(err)
	}

}
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114

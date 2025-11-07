package main

import (
	"log"
	"os"
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"

	"github.com/joho/godotenv"
)

func loadEnv() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable get get work dir : ", err)
	}

	envPath := filepath.Join(cwd, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file")
	} else {
		log.Printf("Loaded .env file")
	}
}
func main() {

	//Load .env file
	loadEnv()
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

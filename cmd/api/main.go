package main

import (
	"log"
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	rootDir := utils.MushGetWorkingDir()
	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	logger.InitLogger(logger.LoggerConfig{
		Level:     "info",
		FileName:  logFile,
		MaxSize:   2,
		MaxBackUp: 5,
		MaxAge:    5,
		Compress:  true,
		IsDev:     utils.GetEnv("APP_ENV", "development"),
	})
	//Load .env file
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Printf("Error loading .env file")
		logger.Log.Warn().Msg("Error loading .env file in api server")
	} else {
		log.Printf("Loaded .env file")
		logger.Log.Info().Msg("Loaded .env file in api server")
	}

	//initialize the config
	cfg := config.NewConfig()
	//initialize theapplicationn
	application, err := app.NewApplication(cfg)
	if err != nil {
		// panic(err)
		logger.Log.Fatal().Err(err).Msgf("Failed to create application: %s", err.Error())
	}

	//start server
	// Run the application
	if err := application.Run(); err != nil {
		// panic(err)
		logger.Log.Fatal().Err(err).Msgf("Application failed to run: %s", err.Error())
	}

}

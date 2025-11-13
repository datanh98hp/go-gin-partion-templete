package main

import (
	"log"
	"os"
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

func mushGetWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable get get work dir : ", err)
	}
	return cwd
}
func loadEnv(path string) {

	err := godotenv.Load(path)
	if err != nil {
		log.Printf("Error loading .env file")
		logger.Log.Warn().Msg("Error loading .env file")
	} else {
		log.Printf("Loaded .env file")
		logger.Log.Info().Msg("Loaded .env file")
	}
}
func main() {
	rootDir := mushGetWorkingDir()
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
	loadEnv(filepath.Join(rootDir, ".env"))
	// Run the application

	//initialize the config
	cfg := config.NewConfig()
	//initialize theapplicationn
	application := app.NewApplication(cfg)

	//start server
	if err := application.Run(); err != nil {
		// panic(err)
		logger.Log.Fatal().Err(err).Msgf("Error: %s", err.Error())
		log.Fatal(err)
	}

}

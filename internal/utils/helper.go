<<<<<<< HEAD
package utils

import (
	"os"
	"user-management-api/pkg/logger"

	"github.com/rs/zerolog"
)

// Helper functions can be added here

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
func NewLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:     level,
		FileName:  path,
		MaxSize:   2,
		MaxBackUp: 5,
		MaxAge:    5,
		Compress:  true,
		IsDev:     GetEnv("APP_ENV", "development"),
	}
	return logger.NewLogger(config)
}
=======
package utils

import (
	"os"
	"strconv"
	"user-management-api/pkg/logger"

	"github.com/rs/zerolog"
)

// Helper functions can be added here

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return valueInt
}
func NewLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:     level,
		FileName:  path,
		MaxSize:   2,
		MaxBackUp: 5,
		MaxAge:    5,
		Compress:  true,
		IsDev:     GetEnv("APP_ENV", "development"),
	}
	return logger.NewLogger(config)
}
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114

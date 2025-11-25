package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
func NewLoggerWithPath(fileName string, level string) *zerolog.Logger {
	// get path to save log
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable get get work dir : ", err)
	}
	path := filepath.Join(cwd, "internal/logs", fileName)

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
func GenerateRandomeString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
func SenitizeRequestBody(data map[string]any, sensitiveKeys []string) map[string]any {
	senitive := make(map[string]any)
	for key, value := range data {
		lowerKey := strings.ToLower(key)
		shoudMask := false
		for _, s := range sensitiveKeys {
			if lowerKey == s {
				shoudMask = true
				break
			}
		}
		if shoudMask {
			senitive[key] = "*****"
		} else {
			switch v := value.(type) {
			case map[string]any:
				senitive[key] = SenitizeRequestBody(v, sensitiveKeys)
			case []any:
				var sentitizeSlice []any
				for _, item := range v {
					if m, ok := item.(map[string]any); ok {
						sentitizeSlice = append(sentitizeSlice, SenitizeRequestBody(m, sensitiveKeys))
					} else {
						sentitizeSlice = append(sentitizeSlice, item)
					}
				}
				senitive[key] = sentitizeSlice
			default:
				senitive[key] = value
			}

		}
	}
	return senitive
}

func MushGetWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable get get work dir : ", err)
	}
	return cwd
}

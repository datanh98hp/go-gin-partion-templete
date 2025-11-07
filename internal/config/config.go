package config

import (
	"fmt"
	"user-management-api/internal/utils"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SslMode  string
}
type Config struct {
	ServerAddress string
	DB            DatabaseConfig
}

func NewConfig() *Config {
	return &Config{
		ServerAddress: utils.GetEnv("PORT", "localhost:8080"),
		DB: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"), // Ensure this is converted to int as needed
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASS", "postgres"),
			DBName:   utils.GetEnv("DB_NAME", "postgres"),
			SslMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
	}
}

func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.SslMode)
}

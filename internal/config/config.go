package config

import (
	"fmt"
	"os"
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
	ServerAddress      string
	DB                 DatabaseConfig
	MailProviderType   string
	MailProviderConfig map[string]any
}

func NewConfig() *Config {
	mailProviderConfig := make(map[string]any)

	mailProviderType := utils.GetEnv("MAIL_PROVIDER_TYPE", "mailtrap")
	if mailProviderType == "mailtrap" {
		mailtrapConfig := map[string]any{
			"mail_sender":      utils.GetEnv("MAILTRAP_MAIL_SENDER", "admin@codewithtuan.com"),
			"name_sender":      utils.GetEnv("MAILTRAP_NAME_SENDER", "Support Team Code With DA"),
			"mailtrap_url":     utils.GetEnv("MAILTRAP_URL", "https://sandbox.api.mailtrap.io/api/send/4201340"),
			"mailtrap_api_key": utils.GetEnv("MAILTRAP_API_KEY", "422fabfb548321a0ab25c49f847badc3"),
		}

		mailProviderConfig["mailtrap"] = mailtrapConfig
	}
	return &Config{
		ServerAddress: fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		DB: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"), // Ensure this is converted to int as needed
			User:     utils.GetEnv("DB_USER", "root"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			DBName:   utils.GetEnv("DB_NAME", "myapp"),
			SslMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
		MailProviderType:   mailProviderType,
		MailProviderConfig: mailProviderConfig,
	}
}

func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.SslMode)
}

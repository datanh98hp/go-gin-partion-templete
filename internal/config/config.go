package config

type Config struct {
	Host          string
	ServerAddress string
}

func NewConfig() *Config {
	return &Config{
		Host:          "localhost",
		ServerAddress: ":8080",
	}
}

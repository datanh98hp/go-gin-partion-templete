package config

import (
	"context"
	"log"
	"time"
	"user-management-api/internal/utils"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Address  string
	Username string
	Password string
	DB       int
}

var client *redis.Client

func NewRedisClient() *redis.Client {
	conf := &RedisConfig{
		Address:  utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		Username: utils.GetEnv("REDIS_USER", ""),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       utils.GetIntEnv("REDIS_DB", 0),
	}

	client := redis.NewClient(&redis.Options{
		Addr:         conf.Address,
		Username:     conf.Address,
		Password:     conf.Password, // no password set
		DB:           conf.DB,       // use default DB
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(context).Result()
	if err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
	}
	log.Printf("Redis connection established successfully")
	return client
}

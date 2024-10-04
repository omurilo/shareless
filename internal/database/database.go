package database

import (
	"context"
	"log"
	"os"

	redis "github.com/redis/go-redis/v9"
)

func NewDbClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
		DB:   0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	return redisClient
}

package database

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	redis "github.com/redis/go-redis/v9"
)

func NewDbClient() *redis.Client {
	addr, db := getHostPortWithDefaults()
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	return redisClient
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getHostPortWithDefaults() (string, int) {
	redisUrl := getEnv("REDIS_URL", "localhost:6379")
	u, err := url.Parse(redisUrl)
	if err != nil {
    return "localhost:6379", 0
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
		port = getEnv("REDIS_PORT", "6379")
	}

	if host == "" {
		host = "localhost"
	}

	redisDB := u.Path
	if redisDB == "" {
		redisDB = fmt.Sprintf("/%s", getEnv("REDIS_DB", "0"))
	}

	db, err := strconv.Atoi(redisDB[1:])
	if err != nil {
		db = 0
	}

  addr := net.JoinHostPort(host, port)

	return addr, db
}

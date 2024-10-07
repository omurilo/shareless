package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	redis "github.com/redis/go-redis/v9"
)

func NewDbClient() *redis.Client {
	redisClient := redis.NewClient(getRedisConfig())

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
	if _, ok := os.LookupEnv("REDIS_URL"); !ok {
		log.Println("REDIS_URL is not present, using localhost:6379")
	}
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

func getRedisConfig() *redis.Options {
	isTLS, err := strconv.ParseBool(os.Getenv("ENABLE_TLS"))
  if err != nil {
    isTLS = false
  }
	addr, db := getHostPortWithDefaults()

	opts := &redis.Options{
		Addr: addr,
		DB:   db,
	}

	if isTLS {
		opts.TLSConfig = &tls.Config{}
	}

	return opts
}

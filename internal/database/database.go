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
	"strings"

	redis "github.com/redis/go-redis/v9"
)

func NewDbClient() *redis.Client {
	redisClient := redis.NewClient(getRedisConfig())

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error connecting to Redis: %+v", err)
	}

	return redisClient
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
  if _, ok := os.LookupEnv("REDIS_URL"); !ok {
		log.Printf("%s is not defined, using %s instead\n", key, defaultValue)
	}
	return defaultValue
}

func getHostPortWithDefaults() (string, string, int) {
	redisUrl := getEnv("REDIS_URL", "redis://localhost:6379")
	if !strings.Contains(redisUrl, "://") {
		redisUrl = "redis://" + redisUrl
	}

	u, err := url.Parse(redisUrl)
	if err != nil {
		return "localhost", "6379", 0
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

	return host, port, db
}

func getRedisConfig() *redis.Options {
	isTLS, err := strconv.ParseBool(os.Getenv("ENABLE_TLS"))
	insecureSkipVerify, err := strconv.ParseBool(os.Getenv("ENABLE_INSECURE_SKIP_VERIFY"))
	if err != nil {
		isTLS = false
	}
	host, port, db := getHostPortWithDefaults()

	opts := &redis.Options{
		Addr: net.JoinHostPort(host, port),
		DB:   db,
	}

	if isTLS {
		opts.TLSConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: insecureSkipVerify,
		}
	}

	return opts
}

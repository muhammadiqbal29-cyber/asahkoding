package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Println("WARNING: REDIS_URL is empty in .env. Skipping Redis connection.")
		return
	}

	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Printf("ERROR: Failed to parse Redis URL: %v\n", err)
		return
	}

	client := redis.NewClient(opts)

	// Ping to test connection with timeout (failover mechanism)
	ctx, cancel := context.WithTimeout(Ctx, 3*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Printf("FAILOVER: Failed to connect to Redis: %v\n", err)
		log.Println("WARNING: API will continue running, but queue features will be disabled.")
		return
	}

	RedisClient = client
	log.Println("SUCCESS: Connected to Redis")
}

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

	// Menambahkan mekanisme Retry agar tangguh terhadap keterlambatan Redis (Docker-in-Docker issue)
	for i := 1; i <= 10; i++ {
		ctx, cancel := context.WithTimeout(Ctx, 3*time.Second)
		_, err = client.Ping(ctx).Result()
		cancel()
		
		if err == nil {
			break
		}
		log.Printf("Menunggu Redis siap (percobaan %d/10): %v\n", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Printf("FAILOVER: Gagal terhubung ke Redis setelah beberapa kali percobaan: %v\n", err)
		log.Println("WARNING: API will continue running, but queue features will be disabled.")
		return
	}

	RedisClient = client
	log.Println("SUCCESS: Connected to Redis")
}

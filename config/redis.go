package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(host string, port string, password string) *redis.Client {
	addr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("❌ Failed to connect to Redis: %v\n", err)
		// We might not want to panic if Redis is optional, but for now let's print error
		// panic(err)
	} else {
		fmt.Println("✅ Connected to Redis successfully")
	}

	return rdb
}

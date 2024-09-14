package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisDB = redis.NewClient(&redis.Options{
	Addr:     EnvValues.RedisAddress,
	Password: "",             
	DB:       0,                
})

func init() {
	// Ping the Redis server to check if it's running
	_, err := RedisDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

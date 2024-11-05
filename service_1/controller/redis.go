package controller

import (
	"boarding-week2/service_1/config"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	cacheDuration = 5 * time.Second
	ctx           = context.Background()
	redisDB       = redis.NewClient(&redis.Options{
		Addr:     config.EnvValues.RedisAddress,
		Password: "",
		DB:       0,
	})
)

func init() {
	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("Could not connect to Redis: ", err)
	}
}

func storeUserInCache(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisDB.Set(ctx, key, data, cacheDuration).Err()
}

func retrieveUserFromCache(key string, result interface{}) (bool, error) {
	data, err := redisDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	redisDB.Expire(ctx, key, cacheDuration)
	return true, json.Unmarshal([]byte(data), result)
}

func deleteUserFromCache(key string) error {
	return redisDB.Del(ctx, key).Err()
}

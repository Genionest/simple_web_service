package config

import (
	"log"

	"example.com/app/global"
	"github.com/go-redis/redis"
)

func InitRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis, got error %v", err)
	}

	global.RedisDB = RedisClient
}

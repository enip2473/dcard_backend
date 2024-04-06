package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(redisAddr string, redisPass string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Redis address
		Password: redisPass, // No password
		DB:       0,         // Default DB
	})

	err := RedisClient.Ping(Ctx).Err()
	return err
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

var Ctx = context.Background()

package db

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisClient redis client
var RedisClient *redis.Client

// NewRedisClient create redis client
func NewRedisClient() error {
	RedisClient = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         "localhost:6379",
		DB:           0,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	_, err := RedisClient.Ping().Result()
	return err
}

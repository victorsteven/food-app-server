package redisdb

import (
	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client

func NewRedisClient(host, port, password string) (*redis.Client, error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return redisClient, nil
}

func NewRedisDB() *redis.Client {
	return redisClient
}

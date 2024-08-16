package redisCaching

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("[redis(caching)]Couldn't connect, Error: %v", err)
	}

}

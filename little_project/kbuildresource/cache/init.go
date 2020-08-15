package cache

import (
	"github.com/go-redis/redis"
	"time"
)

var (
	redisClient *redis.Client
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DialTimeout: time.Minute,
		ReadTimeout: time.Minute,
		WriteTimeout: time.Minute,

	})
}
package redis

import (
	"final/common/consts"
	"os"

	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client {
	addr := os.Getenv(consts.REDIS_ADDR)
	password := os.Getenv(consts.REDIS_PASSWORD)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return client
}

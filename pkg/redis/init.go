package redis

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisbase *redis.Client = nil

// Redis初始化
func Init() error {
	port := os.Getenv("REDIS_PORT")
	host := os.Getenv("REDIS_HOST")
	user := os.Getenv("REDIS_USER")
	password := os.Getenv("REDIS_PASSWORD")
	dsn := fmt.Sprintf("%s:%s", host, port)
	redisbase = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Username: user,
		Password: password,
		DB:       0,
	})
	_, err := redisbase.Ping(redisbase.Context()).Result()
	if err != nil {
		return err
	}
	return nil
}

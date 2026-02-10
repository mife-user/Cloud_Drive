package res

import (
	"fmt"

	"drive/pkg/conf"

	"github.com/go-redis/redis/v8"
)

var redisbase *redis.Client = nil

// Redis初始化
func Init() error {
	config := conf.GetConfig()
	dsn := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)
	redisbase = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	return nil
}

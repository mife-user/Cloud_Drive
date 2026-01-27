package redis

import (
	"github.com/go-redis/redis/v8"
)

// 获取redis连接
func GetRD() *redis.Client {
	if redisbase == nil {
		panic("redis未初始化")
	}
	return redisbase
}

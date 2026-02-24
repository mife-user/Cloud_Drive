package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/cache"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
)

func (r *fileRepo) CheckUserSize(ctx context.Context, userID uint, totalSize int64) (int64, bool) {
	var err error
	var user domain.User
	userKey := fmt.Sprintf("user_size:%d", userID)
	userJSON, err := r.rd.Get(ctx, userKey).Result()
	if err == nil {
		if cache.IsNullValue(userJSON) {
			logger.Error("查询用户失败", logger.U("user_id", userID))
			return 0, false
		}
		if err = exc.ExcJSONToFile(userJSON, &user); err != nil {
			logger.Error("反序列化用户信息失败", logger.C(err))
		} else {
			newSize := user.NowSize + totalSize
			if newSize > user.Size {
				return newSize, false
			}
			return newSize, true
		}
	}
	// 从数据库查询用户
	if err = r.db.Where("id", userID).First(&user).Error; err != nil {
		if err := cache.CacheNullValue(ctx, r.rd, userKey); err != nil {
			logger.Warn("缓存空值失败", logger.C(err))
		}
		logger.Error("查询用户失败", logger.C(err))
		return 0, false
	}
	// 缓存用户信息
	userJSONIn, err := exc.ExcFileToJSON(user)
	if err != nil {
		logger.Error("序列化用户信息失败", logger.C(err))
	} else {
		ttl := cache.UserCacheConfig.RandomTTL()
		if err = r.rd.Set(ctx, userKey, userJSONIn, ttl).Err(); err != nil {
			logger.Warn("缓存用户信息失败", logger.C(err))
		}
	}
	//检查用户额度是否足够
	newSize := user.NowSize + totalSize
	if newSize > user.Size {
		return newSize, false
	}
	return newSize, true
}

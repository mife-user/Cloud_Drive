package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/cache"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"fmt"
)

// GetUserHeadPath 获取用户头像路径
func (r *userRepo) GetUserHeadPath(ctx context.Context, userName string) (string, error) {
	var err error
	var user domain.User
	var header domain.UserHeader
	headerKey := fmt.Sprintf("header:%s", userName)
	headPath, err := r.rd.Get(ctx, headerKey).Result()
	if err == nil {
		if cache.IsNullValue(headPath) {
			logger.Error("查询用户头像失败", logger.S("user_name", userName))
			return "", errorer.New(errorer.ErrUserNotExist)
		}
		return headPath, nil
	}
	// 查询用户ID
	if err = r.db.Where("user_name = ?", userName).First(&user).Error; err != nil {
		if err = cache.CacheNullValue(ctx, r.rd, headerKey); err != nil {
			logger.Warn("缓存空值失败", logger.C(err))
		}
		logger.Error("查询用户失败", logger.S("user_name", userName), logger.C(err))
		return "", err
	}
	// 查询用户头像
	if err = r.db.Where("user_id = ?", user.ID).First(&header).Error; err != nil {
		if err = cache.CacheNullValue(ctx, r.rd, headerKey); err != nil {
			logger.Warn("缓存空值失败", logger.C(err))
		}
		logger.Error("查询用户头像失败", logger.S("user_name", userName), logger.C(err))
		return "", err
	}
	// 缓存用户头像路径
	ttl := cache.UserCacheConfig.RandomTTL()
	if err = r.rd.Set(ctx, headerKey, header.HeaderPath, ttl).Err(); err != nil {
		logger.Warn("缓存用户头像失败", logger.C(err))
	}
	return header.HeaderPath, nil
}

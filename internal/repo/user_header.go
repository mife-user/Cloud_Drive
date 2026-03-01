package repo

import (
	"context"
	"drive/internal/domain"
	"drive/internal/model"
	"drive/pkg/cache"
	"drive/pkg/logger"
	"fmt"
)

// UpdateHeader 更新用户头像
func (r *userRepo) UpdateHeader(ctx context.Context, header *domain.UserHeader) error {
	var err error
	var headerKey = fmt.Sprintf("header:%s", header.Username)
	var headerModel = &model.UserHeader{
		UserName:   header.Username,
		HeaderPath: header.HeaderPath,
	}
	if err = r.db.Where("username = ?", header.Username).Updates(headerModel).Error; err != nil {
		logger.Error("更新用户头像失败", logger.C(err))
		return err
	}
	//更新缓存，使用带随机偏移的缓存策略
	ttl := cache.UserCacheConfig.RandomTTL()
	if err = r.rd.Set(ctx, headerKey, header.HeaderPath, ttl).Err(); err != nil {
		logger.Error("更新用户头像缓存失败", logger.C(err))
		return err
	}
	return nil
}

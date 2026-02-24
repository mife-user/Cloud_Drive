package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/cache"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
)

func (r *fileRepo) ViewFile(ctx context.Context, fileID uint, userID uint) (*domain.File, error) {
	var err error
	var file domain.File
	var fileJSON string

	userKey := fmt.Sprintf("files:%d", userID)
	fileKey := fmt.Sprintf("file:%d", fileID)
	fileJSON, err = r.rd.HGet(ctx, userKey, fileKey).Result()
	if err == nil {
		if cache.IsHashNullValue(fileJSON) {
			logger.Error("查询文件失败", logger.U("file_id", fileID))
			return nil, errorer.New(errorer.ErrFileNotExist)
		}
		if err = exc.ExcJSONToFile(fileJSON, &file); err != nil {
			logger.Error("反序列化文件失败", logger.C(err))
			return nil, err
		}
		if file.DeletedAt.Valid {
			logger.Info("文件已删除", logger.U("file_id", fileID))
			return nil, errorer.New(errorer.ErrFileDeleted)
		}
	} else {
		if err = r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
			if err := cache.CacheHashNullValue(ctx, r.rd, userKey, fileKey); err != nil {
				logger.Warn("缓存空值失败", logger.C(err))
			}
			logger.Error("查询文件失败", logger.C(err))
			return nil, err
		}
		// 缓存文件信息，使用带随机偏移的缓存策略
		if fileJSON, err = exc.ExcFileToJSON(file); err != nil {
			logger.Error("序列化文件信息失败", logger.C(err))
			return nil, err
		}
		if err = r.rd.HSet(ctx, userKey, fileKey, fileJSON).Err(); err != nil {
			logger.Error("缓存文件信息失败", logger.C(err))
			return nil, err
		}
		ttl := cache.FileCacheConfig.RandomTTL()
		if err = r.rd.Expire(ctx, userKey, ttl).Err(); err != nil {
			logger.Error("设置缓存过期时间失败", logger.C(err))
			return nil, err
		}
	}
	logger.Info("查询文件成功", logger.U("file_id", fileID))
	return &file, nil
}

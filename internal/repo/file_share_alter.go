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

// UpdateFilePermissions 更新文件权限
func (r *fileRepo) UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error {
	var err error
	var file domain.File
	var fileJSON string
	// 检查权限参数是否有效
	if permissions != "public" && permissions != "private" {
		logger.Error("无效的权限参数",
			logger.S("permissions", permissions))
		return errorer.New(errorer.ErrInvalidPermissions)
	}
	//缓存检查文件是否存在
	fileKey := fmt.Sprintf("file:%d", fileID)
	userKey := fmt.Sprintf("files:%d", userID)
	fileJSON, err = r.rd.HGet(ctx, userKey, fileKey).Result()
	if err == nil {
		if cache.IsHashNullValue(fileJSON) {
			logger.Error("查询文件失败",
				logger.S("file_id", fmt.Sprintf("%d", fileID)))
			return errorer.New(errorer.ErrFileNotExist)
		}
		if err = exc.ExcJSONToFile(fileJSON, &file); err != nil {
			logger.Error("反序列化文件信息失败", logger.C(err))
			return err
		}
		if file.DeletedAt.Valid {
			logger.Error("文件已被删除", logger.U("file_id", fileID))
			return errorer.New(errorer.ErrFileDeleted)
		}
	} else {
		// 检查文件是否存在
		if err = r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
			if err := cache.CacheHashNullValue(ctx, r.rd, userKey, fileKey); err != nil {
				logger.Warn("缓存空值失败", logger.C(err))
			}
			logger.Error("查询文件失败",
				logger.S("file_id", fmt.Sprintf("%d", fileID)),
				logger.C(err))
			return err
		}
		if fileJSON, err = exc.ExcFileToJSON(file); err != nil {
			logger.Error("序列化文件信息失败", logger.C(err))
			return err
		}
		if err = r.rd.HSet(ctx, userKey, fileKey, fileJSON).Err(); err != nil {
			logger.Error("缓存文件信息失败", logger.C(err))
			return err
		}
	}

	// 检查文件是否存在
	if err = r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		logger.Error("查询文件失败", logger.C(err))
		return err
	}
	// 检查文件所有者或是否为公共文件
	if file.UserID != userID && permissions != "public" {
		logger.Error("非文件所有者，无法更新权限",
			logger.U("file_id", fileID),
			logger.U("user_id", userID))
		return errorer.New(errorer.ErrNotFileOwner)
	}

	file.Permissions = permissions // 更新文件权限
	// 保存文件权限到数据库
	if err = r.db.Model(&domain.File{}).Where("id = ?", fileID).Update("permissions", file.Permissions).Error; err != nil {
		logger.Error("更新文件权限失败", logger.C(err))
		return err
	}
	// 更新缓存中的文件权限
	fileJSON, err = exc.ExcFileToJSON(file)
	if err != nil {
		logger.Error("序列化文件信息失败", logger.C(err))
		return err
	}
	if err = r.rd.HSet(ctx, userKey, fileKey, fileJSON).Err(); err != nil {
		logger.Error("缓存文件信息失败", logger.C(err))
		return err
	}
	// 设置缓存过期时间，使用带随机偏移的缓存策略
	ttl := cache.FileCacheConfig.RandomTTL()
	if err = r.rd.Expire(ctx, userKey, ttl).Err(); err != nil {
		logger.Error("设置缓存过期时间失败", logger.C(err))
		return err
	}
	logger.Info("更新文件权限成功",
		logger.U("file_id", fileID))
	return nil
}

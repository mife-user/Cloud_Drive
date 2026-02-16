package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

// AccessShare 访问文件分享
func (r *fileRepo) AccessShare(ctx context.Context, shareID string, accessKey string) (*domain.File, error) {
	logger.Info("开始访问分享", logger.S("share_id", shareID))
	var file domain.File
	fileShare, err := r.rdShare(ctx, shareID, accessKey)
	if err != nil {
		logger.Error("查询分享记录失败", logger.S("share_id", shareID), logger.C(err))
		return nil, err
	}

	if err := r.db.Where("id = ?", fileShare.FileID).First(&file).Error; err != nil {
		logger.Error("查询文件失败", logger.S("file_id", fmt.Sprintf("%d", fileShare.FileID)), logger.C(err))
		return nil, err
	}

	logger.Info("访问分享成功", logger.S("share_id", shareID), logger.S("file_id", fmt.Sprintf("%d", file.ID)))
	return &file, nil
}
func (r *fileRepo) rdShare(ctx context.Context, shareID string, accessKey string) (*domain.FileShare, error) {
	var fileShare domain.FileShare
	share, err := r.rd.HGetAll(ctx, "share:"+shareID).Result()
	if err != nil {
		// 缓存未命中，从数据库查询
		if err := r.db.Where("share_id = ?", shareID).First(&fileShare).Error; err != nil {
			logger.Error("查询分享记录失败", logger.S("share_id", shareID), logger.C(err))
			return nil, errorer.New(errorer.ErrShareNotExist)
		}

		if fileShare.AccessKey != accessKey {
			logger.Error("访问密钥不匹配", logger.S("share_id", shareID))
			return nil, errorer.New(errorer.ErrInvalidAccessKey)
		}
		if time.Now().Unix() > fileShare.ExpiresAt {
			logger.Error("分享已过期", logger.S("share_id", shareID))
			return nil, errorer.New(errorer.ErrShareExpired)
		}

		// 将数据库结果更新到缓存
		if err := r.rd.HSet(ctx, "share:"+shareID,
			"access_key", fileShare.AccessKey,
			"expires_at", fileShare.ExpiresAt,
			"file_id", fileShare.FileID).Err(); err != nil {
			logger.Warn("缓存分享记录失败", logger.C(err))
		}
		// 设置缓存过期时间
		if err := r.rd.Expire(ctx, "share:"+shareID, 24*time.Hour).Err(); err != nil {
			logger.Warn("设置缓存过期时间失败", logger.C(err))
		}

		return &fileShare, nil
	}

	// 检查缓存键是否存在
	if accessKeyCache, ok := share["access_key"]; !ok || accessKeyCache != accessKey {
		logger.Error("访问密钥不匹配", logger.S("share_id", shareID))
		return nil, errorer.New(errorer.ErrInvalidAccessKey)
	}

	// 检查过期时间
	if expiresAtStr, ok := share["expires_at"]; !ok || expiresAtStr == "" {
		logger.Error("缓存中缺少过期时间", logger.S("share_id", shareID))
		return nil, errorer.New(errorer.ErrShareExpired)
	} else {
		expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
		if err != nil {
			logger.Error("解析过期时间失败", logger.C(err))
			return nil, errorer.New(errorer.ErrShareExpired)
		}
		if time.Now().Unix() > expiresAt {
			logger.Error("分享已过期", logger.S("share_id", shareID))
			return nil, errorer.New(errorer.ErrShareExpired)
		}
	}

	// 解析文件ID
	if fileIDStr, ok := share["file_id"]; !ok || fileIDStr == "" {
		logger.Error("缓存中缺少文件ID", logger.S("share_id", shareID))
		return nil, errorer.New(errorer.ErrShareNotExist)
	} else {
		fileID64, err := strconv.ParseUint(fileIDStr, 10, 64)
		if err != nil {
			logger.Error("解析文件ID失败", logger.C(err))
			return nil, errorer.New(errorer.ErrShareNotExist)
		}
		fileShare.FileID = uint(fileID64)
	}

	return &fileShare, nil
}

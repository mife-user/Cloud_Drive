package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
	"time"
)

// AccessShare 访问文件分享
func (r *fileRepo) AccessShare(ctx context.Context, shareID string, accessKey string) (*domain.File, error) {
	var err error
	// 获取分享记录
	fileShare, err := r.getShareRecord(ctx, shareID, accessKey)
	if err != nil {
		logger.Error("获取分享记录失败", logger.S("share_id", shareID), logger.C(err))
		return nil, err
	}
	// 获取分享文件
	file, err := r.getShareFile(ctx, fileShare)
	if err != nil {
		logger.Error("获取分享文件失败", logger.S("share_id", shareID), logger.C(err))
		return nil, err
	}
	return file, nil
}

// getShareRecord 获取分享记录
func (r *fileRepo) getShareRecord(ctx context.Context, shareID string, accessKey string) (*domain.FileShare, error) {
	var err error
	var fileShare domain.FileShare
	var shareJSON string
	var shareKey = fmt.Sprintf("share:%s", shareID)
	// 从缓存中获取分享记录
	mapcmd := r.rd.Get(ctx, shareKey)
	if err = mapcmd.Err(); err != nil {
		// 缓存中不存在，查询数据库
		if err = r.db.Where("share_id = ?", shareID).First(&fileShare).Error; err != nil {
			logger.Error("查询分享记录失败", logger.S("share_id", shareID), logger.C(err))
			return nil, err
		}
		// 序列化分享记录
		shareJSON, err = exc.ExcFileToJSON(fileShare)
		if err != nil {
			logger.Error("序列化分享记录失败", logger.C(err))
			return nil, err
		}
		// 缓存分享记录，设置过期时间为3小时
		if err = r.rd.Set(ctx, shareKey, shareJSON, 3*time.Hour).Err(); err != nil {
			logger.Error("缓存分享记录失败", logger.C(err))
			return nil, err
		}
	} else {
		shareJSON, err = mapcmd.Result()
		if err != nil {
			logger.Error("获取缓存分享记录失败", logger.C(err))
			return nil, err
		}
		if err = exc.ExcJSONToFile(shareJSON, &fileShare); err != nil {
			logger.Error("反序列化分享记录失败", logger.C(err))
			return nil, err
		}
	}
	// 验证访问密钥
	if fileShare.AccessKey != accessKey {
		logger.Error("访问密钥不匹配", logger.S("share_id", shareID))
		return nil, errorer.New(errorer.ErrInvalidAccessKey)
	}
	// 验证分享有效期
	if time.Now().Unix() > fileShare.ExpiresAt {
		logger.Error("分享已过期", logger.S("share_id", shareID))
		return nil, errorer.New(errorer.ErrShareExpired)
	}
	// 分享记录有效，返回分享记录
	return &fileShare, nil
}

// getShareFile 获取分享文件
func (r *fileRepo) getShareFile(ctx context.Context, share *domain.FileShare) (*domain.File, error) {
	var err error
	var file domain.File
	var fileJSON string
	fileIDSTR := fmt.Sprintf("file:%d", share.FileID)
	userIDSTR := fmt.Sprintf("files:%d", share.UserID)
	mapcmd := r.rd.HGet(ctx, userIDSTR, fileIDSTR)
	if err = mapcmd.Err(); err != nil {
		// 缓存中不存在，查询数据库
		if err = r.db.Where("id = ?", share.FileID).First(&file).Error; err != nil {
			logger.Error("查询文件失败", logger.C(err))
			return nil, err
		}
		fileJSON, err = exc.ExcFileToJSON(file)
		if err != nil {
			logger.Error("序列化文件信息失败", logger.C(err))
			return nil, err
		}
		// 缓存文件信息，设置过期时间为3小时
		if err = r.rd.HSet(ctx, userIDSTR, fileIDSTR, fileJSON).Err(); err != nil {
			logger.Error("缓存文件信息失败", logger.C(err))
			return nil, err
		}
		if err = r.rd.Expire(ctx, userIDSTR, 3*time.Hour).Err(); err != nil {
			logger.Error("设置缓存过期时间失败", logger.C(err))
			return nil, err
		}
	} else {
		fileJSON, err = mapcmd.Result()
		if err != nil {
			logger.Error("获取缓存文件信息失败", logger.C(err))
			return nil, err
		}
		if err = exc.ExcJSONToFile(fileJSON, &file); err != nil {
			logger.Error("反序列化文件信息失败", logger.C(err))
			return nil, err
		}
		if file.DeletedAt.Valid {
			logger.Error("文件已被删除", logger.U("file_id", share.FileID))
			return nil, errorer.New(errorer.ErrFileDeleted)
		}
	}
	return &file, nil
}

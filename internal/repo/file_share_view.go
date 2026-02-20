package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
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
	// 从缓存中获取分享记录
	if err = r.rd.Get(ctx, fmt.Sprintf("share:%s", shareID)).Scan(&fileShare); err != nil {
		// 缓存中不存在，查询数据库
		if err = r.db.Where("share_id = ?", shareID).First(&fileShare).Error; err != nil {
			logger.Error("查询分享记录失败", logger.S("share_id", shareID), logger.C(err))
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
	// 缓存中存在，返回缓存数据
	return &fileShare, nil
}

// getShareFile 获取分享文件
func (r *fileRepo) getShareFile(ctx context.Context, share *domain.FileShare) (*domain.File, error) {
	var err error
	var file domain.File
	fileIDSTR := fmt.Sprintf("file:%d", share.FileID)
	userIDSTR := fmt.Sprintf("files:%d", share.UserID)
	if err = r.rd.HGet(ctx, userIDSTR, fileIDSTR).Scan(&file); err != nil {
		// 缓存中不存在，查询数据库
		if err = r.db.Where("id = ?", share.FileID).First(&file).Error; err != nil {
			logger.Error("查询文件失败", logger.S("file_id", fmt.Sprintf("%d", share.FileID)), logger.C(err))
			return nil, err
		}
	}
	return &file, nil
}

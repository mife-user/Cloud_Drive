package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"drive/pkg/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ShareFile 生成文件分享
func (r *fileRepo) ShareFile(ctx context.Context, fileID uint, userID uint, owner string) (shareID string, accessKey string, err error) {
	logger.Info("开始分享文件", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.S("user_id", fmt.Sprintf("%d", userID)))

	shareID = uuid.New().String()

	accessKey, err = utils.GenerateRandomString(32)
	if err != nil {
		logger.Error("生成访问密钥失败", logger.C(err))
		return "", "", err
	}

	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	fileShare := &domain.FileShare{
		FileID:    fileID,
		ShareID:   shareID,
		AccessKey: accessKey,
		UserID:    userID,
		Owner:     owner,
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(fileShare).Error; err != nil {
		logger.Error("创建分享记录失败", logger.C(err))
		return "", "", err
	}

	logger.Info("分享文件成功", logger.S("share_id", shareID))
	return shareID, accessKey, nil
}

// AccessShare 访问文件分享
func (r *fileRepo) AccessShare(ctx context.Context, shareID string, accessKey string) (*domain.File, error) {
	logger.Info("开始访问分享", logger.S("share_id", shareID))

	var fileShare domain.FileShare
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

	var file domain.File
	if err := r.db.Where("id = ?", fileShare.FileID).First(&file).Error; err != nil {
		logger.Error("查询文件失败", logger.S("file_id", fmt.Sprintf("%d", fileShare.FileID)), logger.C(err))
		return nil, err
	}

	logger.Info("访问分享成功", logger.S("share_id", shareID), logger.S("file_id", fmt.Sprintf("%d", file.ID)))
	return &file, nil
}

// UpdateFilePermissions 更新文件权限
func (r *fileRepo) UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error {
	logger.Info("开始更新文件权限", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.S("user_id", fmt.Sprintf("%d", userID)))

	var file domain.File
	if err := r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		logger.Error("查询文件失败", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.C(err))
		return errorer.New(errorer.ErrFileNotExist)
	}

	if file.UserID != userID {
		logger.Error("非文件所有者，无法更新权限", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.S("user_id", fmt.Sprintf("%d", userID)))
		return errorer.New(errorer.ErrNotFileOwner)
	}

	file.Permissions = permissions
	if err := r.db.Save(&file).Error; err != nil {
		logger.Error("更新文件权限失败", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.C(err))
		return err
	}

	logger.Info("更新文件权限成功", logger.S("file_id", fmt.Sprintf("%d", fileID)))
	return nil
}

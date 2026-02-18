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
	if err := r.rd.HSet(ctx, "share:"+shareID, "access_key", accessKey,
		"expires_at", expiresAt,
		"file_id", fileID).Err(); err != nil {
		logger.Error("缓存分享记录失败", logger.C(err))
		return "", "", err
	}
	// 设置缓存过期时间
	if err := r.rd.Expire(ctx, "share:"+shareID, 24*time.Hour).Err(); err != nil {
		logger.Warn("设置缓存过期时间失败", logger.C(err))
	}
	logger.Info("分享文件成功", logger.S("share_id", shareID))
	return shareID, accessKey, nil
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

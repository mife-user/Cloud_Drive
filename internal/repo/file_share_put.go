package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"drive/pkg/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ShareFile 生成文件分享
func (r *fileRepo) ShareFile(ctx context.Context, fileID uint, userID uint, owner string) (shareID string, accessKey string, err error) {
	logger.Info("开始分享文件",
		logger.S("file_id", fmt.Sprintf("%d", fileID)),
		logger.S("user_id", fmt.Sprintf("%d", userID)))

	shareID = uuid.New().String() // 生成分享ID

	accessKey, err = utils.GenerateRandomString(32) // 生成访问密钥
	if err != nil {
		logger.Error("生成访问密钥失败", logger.C(err))
		return "", "", err
	}

	expiresAt := time.Now().Add(24 * time.Hour).Unix() // 分享有效期为24小时

	fileShare := &domain.FileShare{ // 创建分享记录
		FileID:    fileID,
		ShareID:   shareID,
		AccessKey: accessKey,
		UserID:    userID,
		Owner:     owner,
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(fileShare).Error; err != nil { // 保存分享记录到数据库
		logger.Error("创建分享记录失败", logger.C(err))
		return "", "", err
	}
	shareJSON, err := exc.ExcFileToJSON(fileShare) // 序列化分享记录
	if err != nil {
		logger.Error("序列化分享记录失败", logger.C(err))
		return "", "", err
	}
	r.rd.Set(ctx, fmt.Sprintf("share:%s", shareID), shareJSON, 24*time.Hour) // 缓存分享记录

	logger.Info("分享文件成功", logger.S("share_id", shareID))
	return shareID, accessKey, nil
}

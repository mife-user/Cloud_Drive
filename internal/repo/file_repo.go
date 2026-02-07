package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type fileRepo struct {
	db *gorm.DB
	rd *redis.Client
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB, rd *redis.Client) domain.FileRepo {
	return &fileRepo{db: db, rd: rd}
}

// 文件上传
func (r *fileRepo) UploadFile(ctx context.Context, files []*domain.File) error {
	if len(files) == 0 {
		logger.Error("上传文件失败: 空切片")
		return fmt.Errorf("空切片")
	}
	if err := r.db.Create(files).Error; err != nil {
		logger.Error("上传文件失败", zap.Error(err))
		return err
	}
	logger.Debug("上传文件成功")
	return nil
}

// 查看文件
func (r *fileRepo) ViewFile(ctx context.Context, userID string) ([]domain.File, error) {
	// 从数据库查询文件信息，仅返回文件名、大小、路径、权限和所有者
	var files []domain.File
	if err := r.db.Select("file_name", "size", "path", "permissions", "owner").
		Where("user_id = ?", userID).
		Find(&files).Error; err != nil {
		logger.Error("查询文件失败", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	logger.Debug("查询文件成功")
	return files, nil
}

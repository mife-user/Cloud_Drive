package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type fileRepo struct {
	db *gorm.DB
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB) domain.FileRepo {
	return &fileRepo{db: db}
}

// 文件上传
func (r *fileRepo) UploadFile(ctx context.Context, file *domain.File) error {
	if err := r.db.Create(file).Error; err != nil {
		logger.Error("上传文件失败", zap.Uint("user_id", file.UserID), zap.Error(err))
		return err
	}
	logger.Debug("上传文件成功")
	return nil
}

// 查看文件
func (r *fileRepo) ViewFile(ctx context.Context, userID string) ([]domain.File, error) {
	var files []domain.File
	if err := r.db.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		logger.Error("查询文件失败", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	logger.Debug("查询文件成功")
	return files, nil
}

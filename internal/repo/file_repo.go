package repo

import (
	"drive/internal/domain"

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
func (r *fileRepo) UploadFile(file *domain.File) error {
	return r.db.Create(file).Error
}

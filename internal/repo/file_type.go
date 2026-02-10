package repo

import (
	"drive/internal/domain"

	"github.com/go-redis/redis/v8"
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

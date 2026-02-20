package migrations

import (
	"gorm.io/gorm"
)

// 运行所有迁移
func Run(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &File{}, &FileShare{}, &FileFavorite{})
}

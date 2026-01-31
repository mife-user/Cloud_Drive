package db

import (
	"drive/internal/domain"
	"drive/migrations"
)

// 自动迁移数据库结构
func AutoMigrate() error {
	return GetDB().AutoMigrate(domain.File{}, domain.User{})
}

// 运行数据库迁移
func Migrate() error {
	db := GetDB()
	return migrations.Run(db)
}

// 回滚数据库迁移到指定版本
func Rollback(v int) error {
	db := GetDB()
	return migrations.Rollback(db, v)
}

// MigrationStatus 获取数据库迁移状态
func MigrationStatus() (string, error) {
	db := GetDB()
	return migrations.Status(db)
}

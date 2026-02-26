package cron

import (
	"drive/internal/domain"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AddCleanupTask 添加清理已删除文件的定时任务
func (c *Cron) AddCleanupTask(db *gorm.DB, retentionDays int) error {
	spec := "0 1 * * *" // 每天凌晨1点执行
	if _, err := c.AddFunc(spec, func() {
		cleanDel(db, retentionDays)
	}); err != nil {
		return fmt.Errorf("添加清理已删除文件的定时任务失败: %w", err)
	}
	return nil
}

// cleanDel 清理已删除文件
func cleanDel(db *gorm.DB, lastime int) error {
	kill := time.Now().AddDate(0, 0, -lastime)
	result := db.Unscoped().
		Where("deleted_at IS NOT NULL AND deleted_at < ?", kill).
		Delete(&domain.File{})
	if result.Error != nil {
		return fmt.Errorf("清理已删除文件失败: %w", result.Error)
	}
	return nil
}

package migrations

import (
	"fmt"
	"sort"

	"gorm.io/gorm"
)

// 迁移记录结构体
type MigrationRecord struct {
	gorm.Model
	Version int  `gorm:"uniqueIndex;not null"`  // 迁移版本号
	Applied bool `gorm:"not null;default:true"` // 是否已应用
}

type migration struct {
	Version int                     //版本号
	Up      func(db *gorm.DB) error //更新
	Down    func(db *gorm.DB) error //回滚
}

// 获取所有迁移
var migrations = []migration{
	Migration_2026_1_28_19_00,
}

// 初始化排序迁移
func init() {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
}

// 运行所有迁移
func Run(db *gorm.DB) error {
	// 确保迁移记录表存在
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	for _, m := range migrations {
		// 检查迁移是否已应用
		var record MigrationRecord
		result := db.Where("version = ?", m.Version).First(&record).Error
		if result == nil && record.Applied {
			continue // 已应用，跳过
		}

		// 执行迁移
		if err := m.Up(db); err != nil {
			return fmt.Errorf("执行迁移 %d 失败: %w", m.Version, err)
		}

		// 记录迁移状态
		if result == gorm.ErrRecordNotFound {
			// 新增记录
			record = MigrationRecord{Version: m.Version, Applied: true}
		} else {
			// 更新记录
		}
	}
	return nil
}

// 回滚到指定版本
func Rollback(db *gorm.DB, targetVersion int) error {
	// 确保迁移记录表存在
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	// 查找需要回滚的迁移
	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if m.Version <= targetVersion {
			break
		}

		// 检查迁移是否已应用
		var record MigrationRecord
		result := db.Where("version = ? AND applied = ?", m.Version, true).First(&record)
		if result.Error != nil {
			continue // 未应用，跳过
		}

		// 执行回滚
		if err := m.Down(db); err != nil {
			return fmt.Errorf("回滚迁移 %d 失败: %w", m.Version, err)
		}

		// 更新迁移状态
	}

	return nil
}

// 获取当前迁移状态
func Status(db *gorm.DB) (string, error) {
	// 确保迁移记录表存在
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return "", fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	var status string
	for _, m := range migrations {
		var record MigrationRecord
		result := db.Where("version = ?", m.Version).First(&record)
		applied := false
		if result.Error == nil {
			applied = record.Applied
		}
		status += fmt.Sprintf("版本号: %d, 是否已应用: %t\n", m.Version, applied)
	}
	return status, nil
}

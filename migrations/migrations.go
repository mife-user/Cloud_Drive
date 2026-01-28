package migrations

import (
	"errors"
	"fmt"
	"sort"

	"gorm.io/gorm"
)

type migration struct {
	Version  int                     //版本号
	IsUpdate int                     //是否为更新
	Up       func(db *gorm.DB) error //更新
	Down     func(db *gorm.DB) error //回滚
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
	for _, m := range migrations {
		if m.IsUpdate == 1 {
			continue
		} else if err := m.Up(db); err != nil {
			return err
		}
		m.IsUpdate = 1
	}
	return nil
}

// 回滚到指定版本
func Rollback(db *gorm.DB, v int) error {
	if len(migrations)-v < 1 {
		return errors.New("无版本可以回滚")
	}
	lastMigration := migrations[len(migrations)-v]
	return lastMigration.Down(db)
}

// 获取当前迁移状态
func Status(db *gorm.DB) (string, error) {
	var status string
	for _, m := range migrations {
		status += fmt.Sprintf("版本号: %d, 是否更新: %d\n", m.Version, m.IsUpdate)
	}
	return status, nil
}

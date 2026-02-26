package bootstrap

import (
	"drive/migrations"
	"drive/pkg/db"
)

// initDatabase 初始化数据库
func (a *Application) initDatabase() error {
	if err := db.Init(); err != nil {
		return err
	}
	a.Database = db.GetDB()
	return nil
}

// runMigrations 运行数据库迁移
func (a *Application) runMigrations() error {
	return migrations.Run(a.Database)
}

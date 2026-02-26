package bootstrap

import (
	"context"
	"drive/pkg/logger"
)

// Shutdown 关闭应用
func (a *Application) Shutdown(ctx context.Context) error {
	logger.Info("应用正在关闭...")
	if a.Database != nil {
		sqlDB, err := a.Database.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	if a.Redis != nil {
		a.Redis.Close()
	}
	logger.Info("应用已关闭")
	return nil
}

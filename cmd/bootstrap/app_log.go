package bootstrap

import "drive/pkg/logger"

// initLogger 初始化日志
func (a *Application) initLogger() error {
	return logger.InitLogger(a.Config)
}

package bootstrap

import (
	"fmt"

	"drive/pkg/cron"
	"drive/pkg/logger"
)

func NewApplication() (*Application, error) {
	var err error
	app := &Application{}

	if err = app.loadConfig(); err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	if err = app.initLogger(); err != nil {
		return nil, fmt.Errorf("初始化日志失败: %w", err)
	}

	if err = app.initDatabase(); err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	if err = app.initRedis(); err != nil {
		return nil, fmt.Errorf("初始化Redis失败: %w", err)
	}

	if err = app.runMigrations(); err != nil {
		return nil, fmt.Errorf("运行数据库迁移失败: %w", err)
	}

	if err = app.initRouter(); err != nil {
		return nil, fmt.Errorf("初始化路由失败: %w", err)
	}

	app.initCron()

	return app, nil
}

func (a *Application) Run() error {
	a.printStartupInfo()
	return a.Engine.Run(fmt.Sprintf(":%d", a.Config.Gin.Port))
}

func (a *Application) initCron() {
	a.Cron = cron.NewCron()
	if err := a.Cron.AddCleanupTask(a.Database, 1); err != nil { // 清理1天前的已删除文件
		logger.Error("添加清理已删除文件的定时任务失败", logger.C(err))
	}
	a.Cron.Start()
}

func (a *Application) printStartupInfo() {
	fmt.Println("应用初始化成功！")
	fmt.Printf("配置环境: %s\n", a.Config.Env)
	fmt.Printf("Gin 模式: %s\n", a.Config.Gin.Mode)
	fmt.Printf("服务端口: %d\n", a.Config.Gin.Port)
}

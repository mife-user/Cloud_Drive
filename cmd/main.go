package main

import (
	"fmt"
	"log"

	"drive/internal/api/routes"
	"drive/migrations"
	"drive/pkg/conf"
	"drive/pkg/db"
	"drive/pkg/logger"
)

func main() {
	// 加载配置
	config, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 检查配置
	if err := conf.StatusConfig(); err != nil {
		log.Fatalf("配置检查失败: %v", err)
	}

	// 初始化数据库
	if err := db.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 获取数据库连接
	database := db.GetDB()

	// 运行数据库迁移
	if err := migrations.Run(database); err != nil {
		log.Fatalf("运行数据库迁移失败: %v", err)
	}
	// 初始化日志
	if err := logger.InitLogger(config); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	// 打印迁移状态
	status, err := migrations.Status(database)
	if err != nil {
		log.Printf("获取迁移状态失败: %v", err)
	} else {
		fmt.Println("迁移状态:")
		fmt.Println(status)
	}

	// 初始化路由
	router := routes.GetRouter()
	if !router.NewRouter(database, config) {
		log.Fatalf("创建路由失败")
	}
	engine := router.Setup()

	// 打印初始化成功信息
	fmt.Println("应用初始化成功！")
	fmt.Printf("配置环境: %s\n", config.Env)
	fmt.Printf("Gin 模式: %s\n", config.Gin.Mode)
	fmt.Printf("服务端口: %d\n", config.Gin.Port)

	// 启动服务器
	if err := engine.Run(fmt.Sprintf(":%d", config.Gin.Port)); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

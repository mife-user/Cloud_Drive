package db

import (
	"drive/pkg/conf"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var database *gorm.DB = nil

// 数据库初始化
func Init() error {
	g := conf.GetConfig()
	//重试机制
	for i := 0; i < 5; i++ {
		fmt.Printf("数据库连接尝试 %d\n", i+1)
		conn, err := gorm.Open(mysql.Open(g.Mysql.Dsn), &gorm.Config{})
		if err != nil {
			time.Sleep(8 * time.Second)
			continue
		}
		sqlDB, err := conn.DB()
		if err != nil {
			time.Sleep(8 * time.Second)
			return err
		}
		sqlDB.SetMaxIdleConns(g.Mysql.MaxIdle) //设置空闲连接池中连接的最大数量
		sqlDB.SetMaxOpenConns(g.Mysql.MaxOpen) //设置打开数据库连接的最大数量
		database = conn
		return nil
	}
	return fmt.Errorf("数据库连接失败")
}

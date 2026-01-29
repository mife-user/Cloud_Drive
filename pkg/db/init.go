package db

import (
	"drive/pkg/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var database *gorm.DB

// 数据库初始化
func Init() error {
	g := conf.GetConfig()
	conn, err := gorm.Open(mysql.Open(g.Mysql.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(g.Mysql.MaxIdle) //设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(g.Mysql.MaxOpen) //设置打开数据库连接的最大数量
	database = conn
	return nil
}

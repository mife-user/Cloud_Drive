package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var database *gorm.DB

// 数据库初始化
func Init() error {
	var err error
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)
	database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

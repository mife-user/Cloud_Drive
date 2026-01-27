package db

import "gorm.io/gorm"

//数据库连接
func GetDB() *gorm.DB {
	if database == nil {
		panic("数据库未初始化")
	}
	return database
}

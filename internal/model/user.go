package model

import "gorm.io/gorm"

type UserHeader struct {
	gorm.Model
	UserName   string `gorm:"index"`
	HeaderPath string
}

type User struct {
	gorm.Model
	UserName string
	PassWord string
	Role     string `gorm:"default:NOVIP"`
	Size     int64
	NowSize  int64
}

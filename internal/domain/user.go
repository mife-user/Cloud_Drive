package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string
	PassWord string
	Role     string
}

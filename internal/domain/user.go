package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"username"`
	PassWord string `json:"-"`
	Role     string `json:"role"`
}

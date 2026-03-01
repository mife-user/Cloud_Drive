package domain

import (
	"drive/pkg/errorer"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string
	PassWord string
	Role     string
	Size     int64
	NowSize  int64
}
type UserHeader struct {
	gorm.Model
	UserID     uint
	HeaderPath string
}

func (u *User) IsNullValue() error {
	if u.UserName == "" {
		return errorer.New(errorer.ErrUserNameNotFound)
	}
	if u.PassWord == "" {
		return errorer.New(errorer.ErrPasswordNotFound)
	}
	return nil
}

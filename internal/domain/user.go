package domain

import (
	"drive/pkg/errorer"
)

type User struct {
	ID          uint
	OldUserName string
	UserName    string
	PassWord    string
	Role        string
	Size        int64
	NowSize     int64
}
type UserHeader struct {
	UserName   string
	HeaderPath string
	Role       string
}

// IsNullValue 检查 User 结构体的必填字段是否为空
func (u *User) IsNullValue() error {
	if u.UserName == "" {
		return errorer.New(errorer.ErrUserNameNotFound)
	}
	if u.PassWord == "" {
		return errorer.New(errorer.ErrPasswordNotFound)
	}
	return nil
}

// IsNullValue 检查 UserHeader 结构体的必填字段是否为空
func (u *User) IsNullOldValue() error {
	if u.OldUserName == "" {
		return errorer.New(errorer.ErrOldUserNameNotFound)
	}
	return nil
}

package domain

import (
	"drive/pkg/errorer"
)

type User struct {
	ID       uint
	UserName string
	PassWord string
	Role     string
	Size     int64
	NowSize  int64
}
type UserHeader struct {
	Username   string
	HeaderPath string
	Role       string
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

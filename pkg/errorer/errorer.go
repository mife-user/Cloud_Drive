package errorer

import "errors"

const (
	//用户错误
	ErrUserNameExist    = "用户名已存在"  // ErrUserNameExist 用户名已存在
	ErrPasswordError    = "密码错误"    // ErrPasswordError 密码错误
	ErrUserNotExist     = "用户不存在"   // ErrUserNotExist 用户不存在
	ErrUserNameNotFound = "用户名不可为空" // ErrUserNameNotFound 用户名不可为空
	ErrUpdateUserFailed = "更新用户失败"  // ErrUpdateUserFailed 更新用户失败
	ErrFileNotExist     = "文件不存在"   // ErrFileNotExist 文件不存在
	ErrPasswordNotFound = "密码不可为空"  // ErrPasswordNotFound 密码不可为空
	ErrEmptySlice       = "空切片"     // ErrEmptySlice 空切片
)

func New(err string) error {
	return errors.New(err)
}

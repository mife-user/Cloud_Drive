package errorer

import "errors"

const (
	//用户错误
	ErrUserNameExist    = "用户名已存在"  // ErrUserNameExist 用户名已存在
	ErrPasswordError    = "密码错误"    // ErrPasswordError 密码错误
	ErrUserNotExist     = "用户不存在"   // ErrUserNotExist 用户不存在
	ErrUserNameNotFound = "用户名不可为空" // ErrUserNameNotFound 用户名不可为空
	ErrUpdateUserFailed = "更新用户失败"  // ErrUpdateUserFailed 更新用户失败
	// 文件错误
	ErrFileSizeExceeded = "文件大小超出限制" // ErrFileSizeExceeded 文件大小超出限制
	ErrFileNotExist     = "文件不存在"    // ErrFileNotExist 文件不存在
	ErrPasswordNotFound = "密码不可为空"   // ErrPasswordNotFound 密码不可为空
	ErrEmptySlice       = "空切片"      // ErrEmptySlice 空切片
	// 文件共享错误
	ErrShareNotExist    = "共享不存在"  // ErrShareNotExist 共享不存在
	ErrInvalidAccessKey = "访问密钥无效" // ErrInvalidAccessKey 访问密钥无效
	ErrShareExpired     = "共享已过期"  // ErrShareExpired 共享已过期
	ErrNotFileOwner     = "非文件所有者" // ErrNotFileOwner 非文件所有者
	// 文件收藏错误
	ErrFavoriteExist    = "文件已收藏" // ErrFavoriteExist 文件已收藏
	ErrFavoriteNotExist = "收藏不存在" // ErrFavoriteNotExist 收藏不存在
	// 文件权限错误
	ErrInvalidPermissions = "无效的权限参数" // ErrInvalidPermissions 无效的权限参数
	// 其他错误
	ErrTypeError = "类型错误" // ErrTypeError 类型错误
)

// New 创建一个新的错误
func New(err string) error {
	return errors.New(err)
}

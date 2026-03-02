package request

import "drive/internal/domain"

// UserAlterDT 用户修改DTO
type UserAlterDT struct {
	OldUserName string `json:"old_user_name"`
	UserName    string `json:"user_name"`
	PassWord    string `json:"pass_word"`
}

// 将 UserAlterDT 转换为 domain.User
func (u *UserAlterDT) ToDMUserAlter() *domain.User {
	return &domain.User{
		OldUserName: u.OldUserName,
		UserName:    u.UserName,
		PassWord:    u.PassWord,
	}
}

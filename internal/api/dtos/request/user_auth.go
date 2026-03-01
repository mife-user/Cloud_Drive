package request

import "drive/internal/domain"

// UserAuthDT 用户DTO
type UserAuthDT struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
}

// 将 UserDtos 转换为 domain.User
func (u *UserAuthDT) ToDMUserAuth() *domain.User {
	return &domain.User{
		UserName: u.UserName,
		PassWord: u.PassWord,
	}
}

package request

import "drive/internal/domain"

// UserHeaderDT 用户头像DTO
type UserHeaderDT struct {
	UserName   string
	Role       string
	HeaderPath string
}

// 将 UserHeaderDT 转换为 domain.UserHeader
func (u *UserHeaderDT) ToDMUserHeader() *domain.UserHeader {
	return &domain.UserHeader{
		UserName:   u.UserName,
		Role:       u.Role,
		HeaderPath: u.HeaderPath,
	}
}

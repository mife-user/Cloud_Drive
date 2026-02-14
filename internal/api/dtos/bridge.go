package dtos

import "drive/internal/domain"

// 将 UserDtos 转换为 domain.User
func (u *UserDtos) ToDMUser() *domain.User {
	return &domain.User{
		UserName: u.UserName,
		PassWord: u.PassWord,
	}
}

// 将 FileDtos 转换为 domain.File
func (f *FileDtos) ToDMFile() *domain.File {
	return &domain.File{
		Permissions: f.Permissions,
	}
}

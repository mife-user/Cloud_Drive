package dtos

import "drive/internal/domain"

//将 UserDtos 转换为 domain.User
func (u *UserDtos) ToUser() *domain.User {
	return &domain.User{
		UserName: u.UserName,
		PassWord: u.PassWord,
		Role:     u.Role,
	}
}

//将 FileDtos 转换为 domain.File
func (f *FileDtos) ToFile() *domain.File {
	return &domain.File{
		FileName:    f.FileName,
		Size:        f.Size,
		Path:        f.Path,
		Permissions: f.Permissions,
		Owner:       f.Owner,
	}
}

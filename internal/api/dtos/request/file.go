package request

import "drive/internal/domain"

type FileDtos struct {
	Permissions string `json:"permissions"`
}

type ShareFileRequest struct {
	FileID uint `json:"file_id"`
}

type AccessShareRequest struct {
	AccessKey string `json:"access_key"`
}

type FavoriteFileRequest struct {
	FileID    uint   `json:"file_id"`
	AccessKey string `json:"access_key,omitempty"`
}

// 将 FileDtos 转换为 domain.File
func (f *FileDtos) ToDMFile() *domain.File {
	return &domain.File{
		Permissions: f.Permissions,
	}
}

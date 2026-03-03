package request

import "drive/internal/domain"

// FilePermissionsDT 文件权限DTO
type FilePermissionsDT struct {
	Permissions string `json:"permissions"`
}

// ShareFileDT 分享文件DTO
type ShareFileDT struct {
	FileID uint `json:"file_id"`
}

// AccessShareDT 访问分享DTO
type AccessShareDT struct {
	AccessKey string `json:"access_key"`
}

// FavoriteFileDT 收藏文件DTO
type FavoriteFileDT struct {
	FileID    uint   `json:"file_id"`
	AccessKey string `json:"access_key,omitempty"`
}

// ToDMFilePermissions 转换为domain.File
func (f *FilePermissionsDT) ToDMFilePermissions() *domain.File {
	return &domain.File{
		Permissions: f.Permissions,
	}
}

// ToDMShareFile 转换为domain.FileShare
func (s *ShareFileDT) ToDMShareFile(userID uint, owner string) *domain.FileShare {
	return &domain.FileShare{
		FileID: s.FileID,
		UserID: userID,
		Owner:  owner,
	}
}

// ToDMAccessShare 转换为domain.FileShare
func (a *AccessShareDT) ToDMAccessShare() *domain.FileShare {
	return &domain.FileShare{
		AccessKey: a.AccessKey,
	}
}

// ToDMFavoriteFile 转换为domain.FileFavorite
func (f *FavoriteFileDT) ToDMFavoriteFile(userID uint) *domain.FileFavorite {
	return &domain.FileFavorite{
		UserID: userID,
		FileID: f.FileID,
	}
}

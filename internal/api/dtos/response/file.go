package response

import "drive/internal/domain"

// FileInfoRPS 文件信息响应DTO
type FileInfoRPS struct {
	FileID      uint   `json:"file_id"`
	FileName    string `json:"file_name"`
	Size        int64  `json:"size"`
	Path        string `json:"path"`
	Owner       string `json:"owner"`
	Permissions string `json:"permissions"`
}

// FileShareRPS 文件分享响应DTO
type FileShareRPS struct {
	ShareID   string `json:"share_id"`
	AccessKey string `json:"access_key"`
	FileID    uint   `json:"file_id"`
	Owner     string `json:"owner"`
	ExpiresAt int64  `json:"expires_at"`
}

// FileListRPS 文件列表响应DTO
type FileListRPS struct {
	Files []FileInfoRPS `json:"files"`
	Total int64         `json:"total"`
}

// FavoriteListRPS 收藏列表响应DTO
type FavoriteListRPS struct {
	Favorites []FileInfoRPS `json:"favorites"`
	Total     int64         `json:"total"`
}

// ToDTFileInfo 转换为文件信息响应DTO
func ToDTFileInfo(file *domain.File) *FileInfoRPS {
	return &FileInfoRPS{
		FileID:      file.ID,
		FileName:    file.FileName,
		Size:        file.Size,
		Path:        file.Path,
		Owner:       file.Owner,
		Permissions: file.Permissions,
	}
}

// ToDTFileInfoFromValue 从值类型转换为文件信息响应DTO
func ToDTFileInfoFromValue(file domain.File) FileInfoRPS {
	return FileInfoRPS{
		FileID:      file.ID,
		FileName:    file.FileName,
		Size:        file.Size,
		Path:        file.Path,
		Owner:       file.Owner,
		Permissions: file.Permissions,
	}
}

// ToDTFileList 转换为文件列表响应DTO
func ToDTFileList(files []domain.File) *FileListRPS {
	fileInfos := make([]FileInfoRPS, 0, len(files))
	for _, file := range files {
		fileInfos = append(fileInfos, ToDTFileInfoFromValue(file))
	}
	return &FileListRPS{
		Files: fileInfos,
		Total: int64(len(files)),
	}
}

// ToDTFileShare 转换为文件分享响应DTO
func ToDTFileShare(shareID, accessKey string) *FileShareRPS {
	return &FileShareRPS{
		ShareID:   shareID,
		AccessKey: accessKey,
	}
}

// ToDTFavoriteList 转换为收藏列表响应DTO
func ToDTFavoriteList(files []domain.File) *FavoriteListRPS {
	fileInfos := make([]FileInfoRPS, 0, len(files))
	for _, file := range files {
		fileInfos = append(fileInfos, ToDTFileInfoFromValue(file))
	}
	return &FavoriteListRPS{
		Favorites: fileInfos,
		Total:     int64(len(files)),
	}
}

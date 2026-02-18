package dtos

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

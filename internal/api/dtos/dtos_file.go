package dtos

type FileDtos struct {
	Permissions   string `json:"permissions"`
	IsChunked     bool   `json:"is_chunked"`
	FileMD5       string `json:"file_md5"`
	ChunkIndex    int    `json:"chunk_index"`
	TotalChunks   int    `json:"total_chunks"`
	UploadTaskID  uint   `json:"upload_task_id"`
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

type UploadStatusResponse struct {
	UploadTaskID     uint    `json:"upload_task_id"`
	FileName         string  `json:"file_name"`
	FileSize         int64   `json:"file_size"`
	TotalChunks      int     `json:"total_chunks"`
	CompletedChunks  []int   `json:"completed_chunks"`
	Progress         float64 `json:"progress"`
	Status           int     `json:"status"`
}

package domain

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName    string
	Size        int64
	Path        string
	UserID      uint
	Owner       string
	Permissions string
}

type FileShare struct {
	gorm.Model
	FileID    uint
	ShareID   string
	AccessKey string
	UserID    uint
	Owner     string
	ExpiresAt int64
}

type FileFavorite struct {
	gorm.Model
	UserID uint
	FileID uint
}

type UploadTask struct {
	gorm.Model
	UserID           uint
	FileName         string
	FileSize         int64
	FileMD5          string
	TotalChunks      int
	CompletedChunks  string
	Status           int
}

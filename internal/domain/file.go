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

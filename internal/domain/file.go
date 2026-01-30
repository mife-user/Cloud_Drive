package domain

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName    string `json:"file_name"`
	Size        int64  `json:"size"`
	Path        string `json:"path"`
	UserID      uint   `json:"user_id"`
	Owner       User   `json:"owner"`
	Permissions string `json:"permissions"`
}

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

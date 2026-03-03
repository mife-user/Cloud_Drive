package handlers

import (
	"drive/internal/domain"
	"drive/pkg/conf"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileServicer domain.FileServicer
	config       *conf.Config
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileServicer domain.FileServicer, config *conf.Config) *FileHandler {
	return &FileHandler{
		fileServicer: fileServicer,
		config:   config,
	}
}

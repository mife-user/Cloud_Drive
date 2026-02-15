package handlers

import (
	"drive/internal/domain"
	"drive/pkg/conf"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileRepo domain.FileRepo
	config   *conf.Config
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileRepo domain.FileRepo, config *conf.Config) *FileHandler {
	return &FileHandler{
		fileRepo: fileRepo,
		config:   config,
	}
}

package service

import (
	"drive/internal/domain"
	"drive/pkg/conf"
)

// FileServicer 文件处理器
type fileServicer struct {
	fileRepo domain.FileRepo
	config   *conf.Config
}

// NewFileServicer 创建文件服务
func NewFileServicer(fileRepo domain.FileRepo, config *conf.Config) domain.FileServicer {
	return &fileServicer{
		fileRepo: fileRepo,
		config:   config,
	}
}

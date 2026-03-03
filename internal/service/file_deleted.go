package service

import (
	"context"
	"drive/internal/domain"
)

// GetDeletedFiles 获取已删除文件列表
func (s *fileServicer) GetDeletedFiles(ctx context.Context, userID uint) ([]domain.File, error) {
	var err error
	var files []domain.File
	if files, err = s.fileRepo.GetDeletedFiles(ctx, userID); err != nil {
		return nil, err
	}
	return files, nil
}

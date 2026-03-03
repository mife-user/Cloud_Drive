package service

import (
	"context"
	"drive/internal/domain"
)

// ViewFile 查看文件
func (s *fileServicer) ViewFile(ctx context.Context, fileID uint, userID uint) (*domain.File, error) {
	var err error
	var file *domain.File
	if file, err = s.fileRepo.ViewFile(ctx, fileID, userID); err != nil {
		return nil, err
	}
	return file, nil
}

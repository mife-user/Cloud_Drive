package service

import (
	"context"
	"drive/internal/domain"
)

// ViewFilesNote 查看文件列表
func (s *fileServicer) ViewFilesNote(ctx context.Context, userID uint) ([]domain.File, error) {
	var err error
	var files []domain.File
	if files, err = s.fileRepo.ViewFilesNote(ctx, userID); err != nil {
		return nil, err
	}
	return files, nil
}

package service

import (
	"context"
	"drive/internal/domain"
)

// UploadFile 上传文件
func (s *fileServicer) UploadFile(ctx context.Context, files []*domain.File, nowSize int64) error {
	var err error
	if err = s.fileRepo.UploadFile(ctx, files, nowSize); err != nil {
		return err
	}
	return nil
}

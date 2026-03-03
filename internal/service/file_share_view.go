package service

import (
	"context"
	"drive/internal/domain"
)

// AccessShare 访问分享文件
func (s *fileServicer) AccessShare(ctx context.Context, shareID string, accessKey string) (*domain.File, error) {
	var err error
	var file *domain.File
	if file, err = s.fileRepo.AccessShare(ctx, shareID, accessKey); err != nil {
		return nil, err
	}
	return file, nil
}

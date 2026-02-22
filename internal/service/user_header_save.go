package service

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/save"
	"mime/multipart"
)

func UpdateHeader(ctx context.Context, fileHeader *multipart.FileHeader, userID uint, role string) (*domain.UserHeader, error) {
	var rating int
	// 检查文件大小
	if fileHeader.Size > 1024*1024*10 {
		return nil, errorer.New(errorer.ErrFileSizeExceeded)
	}
	// 检查用户角色是否为VIP用户
	if role != "VIP" {
		rating = 512 * 1024
	} else {
		rating = 1024 * 1024
	}
	// 保存文件
	userID, headerPath, err := save.SaveHeader(ctx, fileHeader, userID, rating)
	if err != nil {
		return nil, err
	}
	return &domain.UserHeader{
		UserID:     userID,
		HeaderPath: headerPath,
	}, nil
}

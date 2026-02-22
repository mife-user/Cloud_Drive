package service

import (
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/save"
	"mime/multipart"
)

func UpdateHeader(fileHeader *multipart.FileHeader, userID uint) (*domain.UserHeader, error) {
	// 检查文件大小
	if fileHeader.Size > 1024*1024*10 {
		return nil, errorer.New(errorer.ErrFileSizeExceeded)
	}
	// 保存文件
	userID, headerPath, err := save.SaveHeader(fileHeader, userID)
	if err != nil {
		return nil, err
	}
	return &domain.UserHeader{
		UserID:     userID,
		HeaderPath: headerPath,
	}, nil
}

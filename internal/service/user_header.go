package service

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/save"
	"mime/multipart"
)

func (s *userServicer) UpdateHeader(ctx context.Context, header *domain.UserHeader, fileHeader *multipart.FileHeader) error {

	var rating int
	// 检查文件大小
	if fileHeader.Size > 1024*1024*10 {
		return errorer.New(errorer.ErrFileSizeExceeded)
	}
	// 检查用户角色是否为VIP用户
	if header.Role != "VIP" {
		rating = 512 * 1024
	} else {
		rating = 1024 * 1024
	}
	// 保存文件
	headerPath, err := save.SaveHeader(ctx, fileHeader, header.Username, rating)
	if err != nil {
		return err
	}
	header.HeaderPath = headerPath
	// 更新数据库
	if err = s.userRepo.UpdateHeader(ctx, header); err != nil {
		return err
	}
	return nil
}

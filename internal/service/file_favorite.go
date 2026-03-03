package service

import (
	"context"
	"drive/internal/domain"
)

// AddFavorite 添加收藏
func (s *fileServicer) AddFavorite(ctx context.Context, userID uint, fileID uint) error {
	var err error
	if err = s.fileRepo.AddFavorite(ctx, userID, fileID); err != nil {
		return err
	}
	return nil
}

// RemoveFavorite 取消收藏
func (s *fileServicer) RemoveFavorite(ctx context.Context, userID uint, fileID uint) error {
	var err error
	if err = s.fileRepo.RemoveFavorite(ctx, userID, fileID); err != nil {
		return err
	}
	return nil
}

// GetFavorites 获取收藏列表
func (s *fileServicer) GetFavorites(ctx context.Context, userID uint) ([]domain.File, error) {
	var err error
	var files []domain.File
	if files, err = s.fileRepo.GetFavorites(ctx, userID); err != nil {
		return nil, err
	}
	return files, nil
}

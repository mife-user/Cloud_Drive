package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"fmt"
)

// AddFavorite 添加文件收藏
func (r *fileRepo) AddFavorite(ctx context.Context, userID uint, fileID uint) error {
	logger.Info("开始添加收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))

	var file domain.File
	if err := r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		logger.Error("查询文件失败", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.C(err))
		return errorer.New(errorer.ErrFileNotExist)
	}

	if file.UserID != userID {
		logger.Error("非文件所有者，无法收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
		return errorer.New(errorer.ErrNotFileOwner)
	}

	var existingFavorite domain.FileFavorite
	if err := r.db.Where("user_id = ? AND file_id = ?", userID, fileID).First(&existingFavorite).Error; err == nil {
		logger.Error("文件已收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
		return errorer.New(errorer.ErrFavoriteExist)
	}

	favorite := &domain.FileFavorite{
		UserID: userID,
		FileID: fileID,
	}

	if err := r.db.Create(favorite).Error; err != nil {
		logger.Error("添加收藏失败", logger.C(err))
		return err
	}

	logger.Info("添加收藏成功", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
	return nil
}

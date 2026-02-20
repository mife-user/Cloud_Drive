package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"fmt"
)

// RemoveFavorite 取消文件收藏
func (r *fileRepo) RemoveFavorite(ctx context.Context, userID uint, fileID uint) error {
	var err error

	result := r.db.Where("user_id = ? AND file_id = ?", userID, fileID).Delete(&domain.FileFavorite{})
	if err = result.Error; err != nil {
		logger.Error("取消收藏失败", logger.C(err))
		return err
	}

	if result.RowsAffected == 0 {
		logger.Error("收藏不存在", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
		return errorer.New(errorer.ErrFavoriteNotExist)
	}

	logger.Info("取消收藏成功", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
	return nil
}

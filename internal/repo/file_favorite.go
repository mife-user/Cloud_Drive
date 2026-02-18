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

// RemoveFavorite 取消文件收藏
func (r *fileRepo) RemoveFavorite(ctx context.Context, userID uint, fileID uint) error {
	logger.Info("开始取消收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))

	result := r.db.Where("user_id = ? AND file_id = ?", userID, fileID).Delete(&domain.FileFavorite{})
	if result.Error != nil {
		logger.Error("取消收藏失败", logger.C(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.Error("收藏不存在", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
		return errorer.New(errorer.ErrFavoriteNotExist)
	}

	logger.Info("取消收藏成功", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
	return nil
}

// GetFavorites 获取收藏列表
func (r *fileRepo) GetFavorites(ctx context.Context, userID uint) ([]domain.File, error) {
	logger.Info("开始获取收藏列表", logger.S("user_id", fmt.Sprintf("%d", userID)))

	var favorites []domain.FileFavorite
	if err := r.db.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		logger.Error("查询收藏列表失败", logger.C(err))
		return nil, err
	}

	var files []domain.File
	for _, favorite := range favorites {
		var file domain.File
		if err := r.db.Where("id = ?", favorite.FileID).First(&file).Error; err != nil {
			logger.Debug("查询文件失败，跳过", logger.S("file_id", fmt.Sprintf("%d", favorite.FileID)), logger.C(err))
			continue
		}

		if file.UserID != userID {
			logger.Debug("用户不再有文件访问权限，跳过", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", file.ID)))
			continue
		}

		files = append(files, file)
	}

	logger.Info("获取收藏列表成功", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("count", fmt.Sprintf("%d", len(files))))
	return files, nil
}

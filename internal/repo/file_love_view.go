package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"
	"fmt"
)

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

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
	// 缓存查询收藏记录
	var rdExist bool
	userKey := fmt.Sprintf("lover:%d", userID)
	fileIDSTR := fmt.Sprintf("file:%d", fileID)
	// 缓存中查询收藏记录
	rdExist, err = r.rd.HExists(ctx, userKey, fileIDSTR).Result()
	if err != nil {
		logger.Error("查询缓存收藏记录失败", logger.C(err))
		return err
	}
	var count int64
	if err = r.db.Model(&domain.FileFavorite{}).Where("user_id = ? AND file_id = ?", userID, fileID).Count(&count).Error; err != nil {
		logger.Error("查询收藏记录失败", logger.C(err))
		return err
	}
	if count == 0 {
		logger.Error("收藏不存在", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
		return errorer.New(errorer.ErrFavoriteNotExist)
	}

	// 从数据库中删除收藏记录
	result := r.db.Where("user_id = ? AND file_id = ?", userID, fileID).Delete(&domain.FileFavorite{})
	if err = result.Error; err != nil {
		logger.Error("取消收藏失败", logger.C(err))
		return err
	}
	// 删除缓存收藏记录
	if rdExist {
		if err = r.rd.HDel(ctx, userKey, fileIDSTR).Err(); err != nil {
			logger.Error("删除缓存收藏记录失败", logger.C(err))
			return err
		}
	}
	logger.Info("取消收藏成功", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
	return nil
}

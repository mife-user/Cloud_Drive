package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
	"time"
)

// AddFavorite 添加文件收藏
func (r *fileRepo) AddFavorite(ctx context.Context, userID uint, fileID uint) error {
	logger.Info("开始添加收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))

	file, err := r.getLoveRecord(ctx, userID, fileID)
	if err != nil {
		logger.Error("查询文件失败", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.C(err))
		return errorer.New(errorer.ErrFileNotExist)
	}

	if file.UserID != userID {
		logger.Error("非文件所有者，无法收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
		return errorer.New(errorer.ErrNotFileOwner)
	}
	if err := r.addLoveRecord(ctx, userID, fileID); err != nil {
		logger.Error("添加收藏失败", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.C(err))
		return err
	}

	logger.Info("添加收藏成功", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
	return nil
}

// getLoveRecord 获取文件收藏记录
func (r *fileRepo) getLoveRecord(ctx context.Context, userID uint, fileID uint) (*domain.File, error) {
	var fileLove domain.File
	userKey := fmt.Sprintf("files:%d", userID)
	fileIDSTR := fmt.Sprintf("file:%d", fileID)
	if fileJSON, err := r.rd.HGet(ctx, userKey, fileIDSTR).Result(); err != nil {
		if err = r.db.Where("id = ?", fileID).First(&fileLove).Error; err != nil {
			logger.Error("查询文件失败", logger.S("file_id", fmt.Sprintf("%d", fileID)), logger.C(err))
			return nil, errorer.New(errorer.ErrFileNotExist)
		}
		var fileJSON string
		if fileJSON, err = exc.ExcFileToJSON(fileLove); err != nil {
			logger.Error("序列化文件信息失败", logger.C(err))
			return nil, err
		}
		if err = r.rd.HSet(ctx, userKey, fileIDSTR, fileJSON).Err(); err != nil {
			logger.Error("缓存文件信息失败", logger.C(err))
			return nil, err
		}
		if err = r.rd.Expire(ctx, userKey, 24*time.Hour).Err(); err != nil {
			logger.Error("设置缓存过期时间失败", logger.C(err))
			return nil, err
		}
	} else {
		if err = exc.ExcJSONToFile(fileJSON, &fileLove); err != nil {
			logger.Error("反序列化文件信息失败", logger.C(err))
			return nil, err
		}
	}
	return &fileLove, nil
}
func (r *fileRepo) addLoveRecord(ctx context.Context, userID uint, fileID uint) error {
	userKey := fmt.Sprintf("loves:%d", userID)
	fileIDSTR := fmt.Sprintf("love:%d", fileID)
	if r.rd.Get(ctx, userKey).Err() == nil {
		logger.Info("文件已收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
		return errorer.New(errorer.ErrFavoriteExist)
	}
	if err := r.db.Where("user_id = ? AND file_id = ?", userID, fileID).First(&domain.FileFavorite{}); err == nil {
		logger.Info("文件已收藏", logger.S("user_id", fmt.Sprintf("%d", userID)), logger.S("file_id", fmt.Sprintf("%d", fileID)))
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
	if err := r.rd.Set(ctx, userKey, fileIDSTR, 24*time.Hour).Err(); err != nil {
		logger.Error("缓存收藏记录失败", logger.C(err))
		return err
	}
	return nil
}

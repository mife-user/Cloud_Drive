package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
)

func (r *fileRepo) ViewFile(ctx context.Context, fileID uint, userID uint) (*domain.File, error) {
	var err error
	var file domain.File
	var fileJSON string

	userKey := fmt.Sprintf("files:%d", userID)
	fileKey := fmt.Sprintf("file:%d", fileID)
	mapCmd := r.rd.HGet(ctx, userKey, fileKey)
	if err = mapCmd.Err(); err != nil {
		if err = r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
			logger.Error("查询文件失败", logger.C(err))
			return nil, err
		}
	} else {
		fileJSON = mapCmd.Val()
		if file.DeletedAt.Valid {
			logger.Info("文件已删除", logger.U("file_id", fileID))
			return nil, errorer.New(errorer.ErrFileDeleted)
		}
		if err = exc.ExcJSONToFile(fileJSON, &file); err != nil {
			logger.Error("反序列化文件失败", logger.C(err))
			return nil, err
		}
	}
	logger.Info("查询文件成功", logger.U("file_id", fileID))
	return &file, nil
}

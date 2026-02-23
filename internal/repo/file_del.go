package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"errors"
	"fmt"
	"time"
)

// DeleteFile 删除文件(软删除)
func (r *fileRepo) DeleteFile(ctx context.Context, userID uint, fileID uint) error {
	var err error
	var file domain.File
	var userKey = fmt.Sprintf("files:%d", userID)
	var fileKey = fmt.Sprintf("file:%d", fileID)
	// 按ID锁
	unlock := r.LockByID(fileID)
	defer unlock()
	// 软删除文件记录
	if err = r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		logger.Error("查询文件失败", logger.U("file_id", fileID), logger.C(err))
		return err
	}
	// 检查文件是否已被删除
	if file.DeletedAt.Valid {
		logger.Error("文件已被删除", logger.U("file_id", fileID))
		return errors.New(errorer.ErrFileDeleted)
	}
	// 软删除文件记录
	if err = r.db.Delete(&file).Error; err != nil {
		logger.Error("删除文件记录失败", logger.U("file_id", fileID), logger.C(err))
		return err
	}
	// 从Redis更新文件元数据
	mapCmd := r.rd.HGet(ctx, userKey, fileKey)
	if err = mapCmd.Err(); err != nil {
		logger.Error("查询文件元数据失败", logger.U("file_id", fileID), logger.C(err))
		return err
	} else {
		fileJSON := mapCmd.Val()
		if err = exc.ExcJSONToFile(fileJSON, &file); err != nil {
			logger.Error("反序列化文件失败", logger.C(err))
			return err
		}
		if file.DeletedAt.Valid {
			logger.Error("文件已被删除", logger.U("file_id", fileID))
			return errors.New(errorer.ErrFileDeleted)
		} else {
			file.DeletedAt.Time = time.Now()
			file.DeletedAt.Valid = true
			fileJSON, err = exc.ExcFileToJSON(file)
			if err != nil {
				logger.Error("序列化文件失败", logger.C(err))
				return err
			}
			if err = r.rd.HSet(ctx, userKey, fileKey, fileJSON).Err(); err != nil {
				logger.Error("更新删除缓存数据失败", logger.U("file_id", fileID), logger.C(err))
				return err
			}
		}
	}
	logger.Info("删除文件成功", logger.U("file_id", fileID))
	return nil
}

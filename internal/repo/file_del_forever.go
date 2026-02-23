package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"
	"fmt"
)

// DeleteFileForever 删除文件永久删除
func (r *fileRepo) DeleteFileForever(ctx context.Context, userID uint, fileID uint) error {
	var err error
	var userKey = fmt.Sprintf("files:%d", userID)
	var fileKey = fmt.Sprintf("file:%d", fileID)
	// 从缓存中删除文件
	mapCmd := r.rd.HDel(ctx, userKey, fileKey)
	if err = mapCmd.Err(); err != nil {
		logger.Error("删除文件缓存失败", logger.C(err))
		return err
	} else {
		if err = r.db.Unscoped().Where("id = ?", fileID).Delete(&domain.File{}).Error; err != nil {
			logger.Error("删除文件记录失败", logger.C(err))
			return err
		}
	}
	return nil
}

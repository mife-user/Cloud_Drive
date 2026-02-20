package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"fmt"
)

// UpdateFilePermissions 更新文件权限
func (r *fileRepo) UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error {
	var err error
	var file domain.File

	if err = r.db.Where("id = ?", fileID).First(&file).Error; err != nil { // 查询文件是否存在
		logger.Error("查询文件失败",
			logger.S("file_id", fmt.Sprintf("%d", fileID)),
			logger.C(err))
		return err
	}

	if file.UserID != userID {
		logger.Error("非文件所有者，无法更新权限",
			logger.S("file_id", fmt.Sprintf("%d", fileID)),
			logger.S("user_id", fmt.Sprintf("%d", userID)))
		return errorer.New(errorer.ErrNotFileOwner)
	}

	file.Permissions = permissions // 更新文件权限

	if err = r.db.Save(&file).Error; err != nil { // 保存文件权限到数据库
		logger.Error("更新文件权限失败",
			logger.S("file_id", fmt.Sprintf("%d", fileID)),
			logger.C(err))
		return err
	}

	logger.Info("更新文件权限成功",
		logger.S("file_id", fmt.Sprintf("%d", fileID)))
	return nil
}

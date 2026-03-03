package service

import (
	"context"
)

// UpdateFilePermissions 更新文件权限
func (s *fileServicer) UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error {
	var err error
	if err = s.fileRepo.UpdateFilePermissions(ctx, fileID, userID, permissions); err != nil {
		return err
	}
	return nil
}

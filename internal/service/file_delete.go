package service

import (
	"context"
)

// DeleteFile 删除文件(软删除)
func (s *fileServicer) DeleteFile(ctx context.Context, userID uint, fileID uint) error {
	var err error
	if err = s.fileRepo.DeleteFile(ctx, userID, fileID); err != nil {
		return err
	}
	return nil
}

// DeleteFileForever 永久删除文件
func (s *fileServicer) DeleteFileForever(ctx context.Context, userID uint, fileID uint) error {
	var err error
	if err = s.fileRepo.DeleteFileForever(ctx, userID, fileID); err != nil {
		return err
	}
	return nil
}

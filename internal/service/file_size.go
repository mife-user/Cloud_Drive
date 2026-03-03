package service

import (
	"context"
)

// CheckUserSize 检查用户空间大小
func (s *fileServicer) CheckUserSize(ctx context.Context, userID uint, totalSize int64) (int64, bool) {
	var newSize int64
	var ok bool
	newSize, ok = s.fileRepo.CheckUserSize(ctx, userID, totalSize)
	return newSize, ok
}

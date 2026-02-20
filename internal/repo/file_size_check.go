package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"
)

func (r *fileRepo) CheckUserSize(ctx context.Context, userID uint, totalSize int64) (int64, bool) {
	var user domain.User
	if err := r.db.Where("user_id", userID).First(&user).Error; err != nil {
		logger.Error("查询用户失败", logger.C(err))
		return 0, false
	}
	//检查用户额度是否足够
	newSize := user.NowSize + totalSize
	if newSize > user.Size {
		return newSize, false
	}
	return newSize, true
}

package service

import (
	"context"
	"drive/internal/domain"
)

// RemixUser 修改用户信息
func (s *userServicer) RemixUser(ctx context.Context, user *domain.User) error {
	var err error
	// 检查用户信息是否完整
	err = user.IsNullOldValue()
	if err != nil {
		return err
	}
	// 从上下文获取当前登录用户ID
	if err = s.userRepo.RemixUser(ctx, user); err != nil {
		return err
	}
	return nil
}

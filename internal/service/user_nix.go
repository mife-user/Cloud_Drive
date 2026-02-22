package service

import (
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
)

// 完善用户信息
func BuildUser(DMUser *domain.User, userID any) error {
	// 检查userID是否为uint类型
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		logger.Error(errorer.ErrTypeError)
		return errorer.New(errorer.ErrTypeError)
	}
	DMUser.ID = userIDUint
	return nil
}

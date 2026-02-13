package service

import (
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
)

// 完善用户信息
func BuildUser(DMUser *domain.User, userID any) error {
	var ok bool
	DMUser.ID, ok = userID.(uint)
	if !ok {
		logger.Error(errorer.ErrTypeError)
		return errorer.New(errorer.ErrTypeError)
	}
	return nil
}

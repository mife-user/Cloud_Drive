package service

import "drive/internal/domain"

// 完善用户信息
func BuildUser(DMUser *domain.User, userID any) error {
	DMUser.ID = userID.(uint)
	return nil
}

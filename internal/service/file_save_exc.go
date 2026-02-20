package service

import (
	"drive/pkg/exc"
	"drive/pkg/logger"
)

func ExchangeType(userID any, userName any, userRole any) (uint, string, string, bool) {
	// 转换userID为uint类型
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		logger.Error("userID类型转换失败")
		return 0, "", "", false
	}
	// 转换userName为string类型
	userNameStr, ok := exc.IsString(userName)
	if !ok {
		logger.Error("userName类型转换失败")
		return 0, "", "", false
	}
	// 转换userRole为string类型
	userRoleStr, ok := exc.IsString(userRole)
	if !ok {
		logger.Error("userRole类型转换失败")
		return 0, "", "", false
	}
	return userIDUint, userNameStr, userRoleStr, true
}

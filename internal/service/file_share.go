package service

import (
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
)

// ExchangeFile 兑换文件
func ExchangeFile(userID any, userName any) (uint, string, error) {
	// 转换userID为uint类型
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		logger.Error("userID类型转换失败")
		return 0, "", errorer.New(errorer.ErrTypeError)
	}
	// 转换userName为string类型
	userNameStr, ok := exc.IsString(userName)
	if !ok {
		logger.Error("userName类型转换失败")
		return 0, "", errorer.New(errorer.ErrTypeError)
	}
	return userIDUint, userNameStr, nil
}

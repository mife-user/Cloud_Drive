package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"drive/pkg/utils"
	"time"
)

// 用户注册
func (r *userRepo) Register(ctx context.Context, user *domain.User) error {
	// 检查用户名是否为空
	if user.UserName == "" {
		logger.Debug("注册用户失败"+errorer.ErrUserNameNotFound, logger.S("user_name", user.UserName))
		return errorer.New(errorer.ErrUserNameNotFound)
	}
	// 检查密码是否为空
	if user.PassWord == "" {
		logger.Debug("注册用户失败"+errorer.ErrPasswordNotFound, logger.S("user_name", user.UserName))
		return errorer.New(errorer.ErrPasswordNotFound)
	}
	//缓存检查用户是否已存在
	if err := r.rd.Get(ctx, "user:"+user.UserName).Err(); err == nil {
		logger.Debug("注册用户失败"+errorer.ErrUserNameExist, logger.S("user_name", user.UserName))
		return errorer.New(errorer.ErrUserNameExist)
	}
	// 检查用户名是否已存在
	var existingUser domain.User
	if err := r.db.Where("user_name = ?", user.UserName).First(&existingUser).Error; err == nil {
		// 用户已存在
		logger.Debug("注册用户失败"+errorer.ErrUserNameExist, logger.S("user_name", user.UserName))
		return errorer.New(errorer.ErrUserNameExist)
	}
	// 加密密码
	hashedPassword, err := utils.HashPassword(user.PassWord)
	if err != nil {
		logger.Error("注册用户失败", logger.S("user_name", user.UserName), logger.C(err))
		return err
	}
	user.PassWord = hashedPassword
	// 创建用户
	if err := r.db.Create(user).Error; err != nil {
		logger.Error("注册用户失败", logger.S("user_name", user.UserName), logger.C(err))
		return err
	}
	// 缓存用户所有信息
	userjson, errjson := exc.ExcFileToJSON(user)
	if errjson != nil {
		logger.Error("注册用户失败", logger.S("user_name", user.UserName), logger.C(errjson))
		return errjson
	}
	if err := r.rd.Set(ctx, "user:"+user.UserName, userjson, time.Hour*3).Err(); err != nil {
		logger.Error("缓存用户信息失败", logger.S("user_name", user.UserName), logger.C(err))
		// 缓存失败不影响注册结果
	}
	logger.Info("注册用户成功", logger.S("user_name", user.UserName))
	return nil
}

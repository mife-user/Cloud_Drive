package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
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
	if err := r.db.Where("user_name = ?", user.UserName).First(&domain.User{}).Error; err != nil {
		// 加密密码
		hashedPassword, err := utils.HashPassword(user.PassWord)
		if err != nil {
			logger.Debug("注册用户失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
		user.PassWord = hashedPassword
		// 创建用户
		if err := r.db.Create(user).Error; err != nil {
			logger.Error("注册用户失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
	}

	// 缓存用户信息
	if err := r.rd.Set(ctx, "user:"+user.UserName, user, time.Hour*3).Err(); err != nil {
		logger.Debug("注册用户失败", logger.S("user_name", user.UserName), logger.C(err))
		return err
	}
	logger.Info("注册用户成功", logger.S("user_name", user.UserName))
	return nil
}

// 用户登录
func (r *userRepo) Logon(ctx context.Context, user *domain.User) error {
	// 先根据用户名查询用户
	var existingUser domain.User
	if err := r.rd.Get(ctx, "user:"+user.UserName).Scan(&existingUser); err != nil {
		if err := r.db.Where("user_name = ?", user.UserName).First(&existingUser).Error; err != nil {
			logger.Error("登录用户失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
		// 缓存用户信息
		if err := r.rd.Set(ctx, "user:"+user.UserName, &existingUser, time.Hour*3).Err(); err != nil {
			logger.Debug("登录用户失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
	}
	// 验证密码
	if !utils.CheckPasswordHash(user.PassWord, existingUser.PassWord) {
		logger.Debug("登录用户失败"+errorer.ErrPasswordError, logger.S("user_name", user.UserName))
		return errorer.New(errorer.ErrPasswordError)
	}

	// 将查询到的用户信息赋值给传入的user
	*user = existingUser
	logger.Info("登录用户成功", logger.S("user_name", user.UserName))
	return nil
}

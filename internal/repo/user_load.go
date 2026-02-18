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

// 用户登录
func (r *userRepo) Logon(ctx context.Context, user *domain.User) error {
	// 检查用户名是否为空
	if user.UserName == "" {
		logger.Debug("登录用户失败"+errorer.ErrUserNameNotFound, logger.S("user_name", user.UserName))
		return errorer.New(errorer.ErrUserNameNotFound)
	}
	// 检查密码是否为空
	if user.PassWord == "" {
		logger.Debug("登录用户失败"+errorer.ErrPasswordNotFound, logger.S("user_name", user.UserName))
		return errorer.New(errorer.ErrPasswordNotFound)
	}
	// 先根据用户名查询用户
	var existingUser domain.User
	userjsonOut, errjsonout := r.rd.Get(ctx, "user:"+user.UserName).Result()
	if errjsonout == nil {
		if err := exc.ExcJSONToFile(userjsonOut, &existingUser); err != nil {
			logger.Error("从缓存中解析用户信息失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
	} else {
		// 缓存中不存在用户，从数据库查询
		if err := r.db.Where("user_name = ?", user.UserName).First(&existingUser).Error; err != nil {
			logger.Error("登录用户失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
		// 缓存用户信息
		userjsonIn, errjsonIn := exc.ExcFileToJSON(existingUser)
		if errjsonIn != nil {
			logger.Error("缓存用户信息失败", logger.S("user_name", user.UserName), logger.C(errjsonIn))
			// 缓存失败不影响登录结果
		}
		if err := r.rd.Set(ctx, "user:"+user.UserName, userjsonIn, time.Hour*3).Err(); err != nil {
			logger.Error("缓存用户信息失败", logger.S("user_name", user.UserName), logger.C(err))
			// 缓存失败不影响登录结果
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

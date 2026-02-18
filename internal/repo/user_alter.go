package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"drive/pkg/utils"
	"encoding/json"
	"time"
)

// RemixUser 修改用户信息
func (r *userRepo) RemixUser(ctx context.Context, user *domain.User) error {
	// 检查用户名是否为空
	if user.UserName == "" {
		logger.Debug("修改用户失败" + errorer.ErrUserNameNotFound)
		return errorer.New(errorer.ErrUserNameNotFound)
	}
	// 检查密码是否为空
	if user.PassWord == "" {
		logger.Debug("修改用户失败" + errorer.ErrPasswordNotFound)
		return errorer.New(errorer.ErrPasswordNotFound)
	}
	var oldUser domain.User
	// 检查用户是否存在
	if err := r.db.Where("id = ?", user.ID).First(&oldUser).Error; err == nil {
		// 加密密码
		hashedPassword, err := utils.HashPassword(user.PassWord)
		if err != nil {
			logger.Debug("修改用户失败"+errorer.ErrPasswordError, logger.C(err))
			return err
		}
		user.PassWord = hashedPassword
		//避免修改敏感信息
		user.Role = oldUser.Role
		user.ID = oldUser.ID
		// 更新用户信息
		if err := r.db.Model(&domain.User{}).
			Where("id = ?", user.ID).
			Updates(user).Error; err != nil {
			logger.Error("修改用户失败"+errorer.ErrUpdateUserFailed, logger.C(err))
			return err
		}
		// 缓存用户信息
		userjsonIn, errjson := json.Marshal(user)
		if errjson != nil {
			logger.Error("修改用户失败"+errorer.ErrUpdateUserFailed, logger.C(errjson))
			return errjson
		}
		if err := r.rd.Set(ctx, "user:"+user.UserName, string(userjsonIn), time.Hour*3).Err(); err != nil {
			logger.Debug("修改用户失败"+errorer.ErrUpdateUserFailed, logger.C(err))
			return err
		}
		logger.Info("修改用户成功", logger.S("user_name", user.UserName))
		return nil
	}
	logger.Debug("修改用户失败" + errorer.ErrUserNotExist)
	return errorer.New(errorer.ErrUserNotExist)
}

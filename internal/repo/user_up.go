package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"drive/pkg/utils"
	"errors"
	"time"

	"go.uber.org/zap"
)

func (r *userRepo) RemixUser(ctx context.Context, user *domain.User) error {
	// 检查用户名是否为空
	if user.UserName == "" {
		logger.Error("修改用户失败", zap.Error(errors.New(errorer.ErrUserNameNotFound)))
		return errors.New(errorer.ErrUserNameNotFound)
	}
	// 检查密码是否为空
	if user.PassWord == "" {
		logger.Error("修改用户失败", zap.Error(errors.New(errorer.ErrPasswordNotFound)))
		return errors.New(errorer.ErrPasswordNotFound)
	}
	var oldUser domain.User
	// 检查用户是否存在
	if err := r.db.Where("id = ?", user.ID).First(&oldUser).Error; err == nil {
		// 加密密码
		hashedPassword, err := utils.HashPassword(user.PassWord)
		if err != nil {
			logger.Error("修改用户失败", zap.Error(err))
			return errors.New(errorer.ErrPasswordError)
		}
		user.PassWord = hashedPassword
		//避免修改敏感信息
		user.Role = oldUser.Role
		user.ID = oldUser.ID
		// 更新用户信息
		if err := r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
			logger.Error("修改用户失败", zap.Error(err))
			return errors.New(errorer.ErrUpdateUserFailed)
		}
		// 缓存用户信息
		if err := r.rd.Set(ctx, "user:"+user.UserName, user, time.Hour*3).Err(); err != nil {
			logger.Error("修改用户失败", zap.Error(err))
			return err
		}
		logger.Info("修改用户成功", zap.String("user_name", user.UserName))
		return nil
	}
	logger.Error("修改用户失败", zap.Error(errors.New(errorer.ErrUserNotExist)))
	return errors.New(errorer.ErrUserNotExist)
}

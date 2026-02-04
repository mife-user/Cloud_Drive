package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"
	"drive/pkg/utils"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) domain.UserRepo {
	return &userRepo{db: db}
}

// 用户注册
func (r *userRepo) Register(ctx context.Context, user *domain.User) error {
	if user.UserName == "" || user.PassWord == "" {
		logger.Error("注册用户失败", zap.String("user_name", user.UserName), zap.Error(errors.New("用户名或密码不能为空")))
		return errors.New("用户名或密码不能为空")
	}
	// 加密密码
	hashedPassword, err := utils.HashPassword(user.PassWord)
	if err != nil {
		logger.Error("注册用户失败", zap.String("user_name", user.UserName), zap.Error(err))
		return err
	}
	user.PassWord = hashedPassword

	if err := r.db.Create(user).Error; err != nil {
		logger.Error("注册用户失败", zap.String("user_name", user.UserName), zap.Error(err))
		return err
	}
	logger.Debug("注册用户成功")
	return nil
}

// 用户登录
func (r *userRepo) Logon(ctx context.Context, user *domain.User) error {
	// 先根据用户名查询用户
	var existingUser domain.User
	if err := r.db.Where("user_name = ?", user.UserName).First(&existingUser).Error; err != nil {
		logger.Error("登录用户失败", zap.String("user_name", user.UserName), zap.Error(err))
		return err
	}

	// 验证密码
	if !utils.CheckPasswordHash(user.PassWord, existingUser.PassWord) {
		logger.Error("登录用户失败", zap.String("user_name", user.UserName), zap.Error(errors.New("密码错误")))
		return errors.New("密码错误")
	}

	// 将查询到的用户信息赋值给传入的user
	*user = existingUser
	logger.Debug("登录用户成功")
	return nil
}

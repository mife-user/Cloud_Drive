package repo

import (
	"context"
	"drive/internal/domain"
	"drive/internal/model"
	"drive/pkg/cache"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"drive/pkg/utils"
	"fmt"
)

// 用户注册
func (r *userRepo) Register(ctx context.Context, user *domain.User) error {
	var err error
	//缓存检查用户是否已存在
	existingUser := &model.User{
		UserName: user.UserName,
		PassWord: user.PassWord,
	}
	key := "user:" + existingUser.UserName
	value, err := r.rd.Get(ctx, key).Result()
	if err == nil {
		if cache.IsNullValue(value) {
			if err = r.rd.Del(ctx, key).Err(); err != nil {
				logger.Warn("删除空值缓存失败", logger.C(err))
			}
		} else {
			logger.Debug(fmt.Sprintf("注册用户失败%s", errorer.ErrUserNameExist), logger.S("user_name", existingUser.UserName), logger.C(err))
			return errorer.New(errorer.ErrUserNameExist)
		}
	}
	// 检查用户名是否已存在
	var count int64
	if err = r.db.Model(&model.User{}).Where("user_name = ?", existingUser.UserName).Count(&count).Error; err != nil {
		logger.Error("检查用户是否存在失败", logger.S("user_name", existingUser.UserName), logger.C(err))
		return err
	}
	if count > 0 {
		// 用户已存在
		logger.Debug(fmt.Sprintf("注册用户失败%s", errorer.ErrUserNameExist), logger.S("user_name", existingUser.UserName))
		return errorer.New(errorer.ErrUserNameExist)
	}
	// 加密密码
	hashedPassword, err := utils.HashPassword(existingUser.PassWord)
	if err != nil {
		logger.Error("注册用户失败", logger.S("user_name", existingUser.UserName), logger.C(err))
		return err
	}
	existingUser.PassWord = hashedPassword
	// 创建用户
	if err = r.db.Select("user_name", "pass_word", "role").Create(existingUser).Error; err != nil {
		logger.Error("注册用户失败", logger.S("user_name", existingUser.UserName), logger.C(err))
		return err
	}
	// 缓存用户所有信息
	userjson, errjson := exc.ExcFileToJSON(existingUser)
	if errjson != nil {
		logger.Error("注册用户失败", logger.S("user_name", existingUser.UserName), logger.C(errjson))
		return errjson
	}
	// 使用带随机偏移的缓存策略
	ttl := cache.UserCacheConfig.RandomTTL()
	if err = r.rd.Set(ctx, "user:"+existingUser.UserName, userjson, ttl).Err(); err != nil {
		logger.Error("缓存用户信息失败", logger.C(err))
		// 缓存失败不影响注册结果
	}
	user.ID = existingUser.ID
	user.Role = existingUser.Role
	user.Size = existingUser.Size
	user.NowSize = existingUser.NowSize
	logger.Info("注册用户成功", logger.S("user_name", existingUser.UserName))
	return nil
}

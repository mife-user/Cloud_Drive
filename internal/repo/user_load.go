package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/cache"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"drive/pkg/utils"
)

// 用户登录
func (r *userRepo) Logon(ctx context.Context, user *domain.User) error {
	var err error
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
	key := "user:" + user.UserName
	userjsonOut, err := r.rd.Get(ctx, key).Result()
	if err == nil {
		if cache.IsNullValue(userjsonOut) {
			logger.Debug("登录用户失败"+errorer.ErrUserNotExist, logger.S("user_name", user.UserName))
			return errorer.New(errorer.ErrUserNotExist)
		}
		if err = exc.ExcJSONToFile(userjsonOut, &existingUser); err != nil {
			logger.Error("从缓存中解析用户信息失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
	} else {
		// 缓存中不存在用户，从数据库查询
		if err = r.db.Where("user_name = ?", user.UserName).First(&existingUser).Error; err != nil {
			if err = cache.CacheNullValue(ctx, r.rd, key); err != nil {
				logger.Warn("缓存空值失败", logger.C(err))
			}
			logger.Error("登录用户失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
		// 使用带随机偏移的缓存策略
		userjsonIn, err := exc.ExcFileToJSON(existingUser)
		if err != nil {
			logger.Error("缓存用户信息失败", logger.S("user_name", user.UserName), logger.C(err))
			// 缓存失败不影响登录结果
		}
		ttl := cache.UserCacheConfig.RandomTTL()
		if err = r.rd.Set(ctx, key, userjsonIn, ttl).Err(); err != nil {
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

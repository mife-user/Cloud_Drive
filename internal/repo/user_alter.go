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
)

// RemixUser 修改用户信息
func (r *userRepo) RemixUser(ctx context.Context, user *domain.User) error {
	var err error
	oldUser := &model.User{
		UserName: user.OldUserName,
	}
	newUser := &model.User{
		UserName: user.UserName,
		PassWord: user.PassWord,
	}
	key := "user:" + oldUser.UserName
	// 检查用户是否存在
	userjsonOut, err := r.rd.Get(ctx, key).Result()
	if err == nil {
		if cache.IsNullValue(userjsonOut) {
			logger.Debug("修改用户失败" + errorer.ErrUserNotExist)
			return errorer.New(errorer.ErrUserNotExist)
		}
		if err = exc.ExcJSONToFile(userjsonOut, &oldUser); err != nil {
			logger.Error("从缓存中解析用户信息失败", logger.S("user_name", user.UserName), logger.C(err))
			return err
		}
	} else {
		if err = r.db.Where("id = ?", user.ID).First(&oldUser).Error; err != nil {
			if err := cache.CacheNullValue(ctx, r.rd, key); err != nil {
				logger.Warn("缓存空值失败", logger.C(err))
			}
			logger.Debug("修改用户失败" + errorer.ErrUserNotExist)
			return errorer.New(errorer.ErrUserNotExist)
		}
	}
	// 加密密码
	hashedPassword, err := utils.HashPassword(newUser.PassWord)
	if err != nil {
		logger.Debug("修改用户失败"+errorer.ErrPasswordError, logger.C(err))
		return err
	}
	newUser.PassWord = hashedPassword
	//避免修改敏感信息(虽然不可能修改，但以防万一,hhh)
	newUser.Role = oldUser.Role
	newUser.ID = oldUser.ID
	// 更新用户信息
	if err = r.db.Model(&model.User{}).
		Where("id = ?", newUser.ID).
		Select("user_name", "pass_word").
		Updates(newUser).Error; err != nil {
		logger.Error("修改用户失败"+errorer.ErrUpdateUserFailed, logger.C(err))
		return err
	}
	// 缓存用户信息，使用带随机偏移的缓存策略
	userjsonIn, err := exc.ExcFileToJSON(newUser)
	if err != nil {
		logger.Error("修改用户失败"+errorer.ErrUpdateUserFailed, logger.C(err))
		return err
	}
	ttl := cache.UserCacheConfig.RandomTTL()
	if err = r.rd.Set(ctx, key, userjsonIn, ttl).Err(); err != nil {
		logger.Debug("修改用户失败"+errorer.ErrUpdateUserFailed, logger.C(err))
		return err
	}
	logger.Info("修改用户成功", logger.S("user_name", user.UserName))
	return nil
}

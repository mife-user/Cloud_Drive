package repo

import (
	"drive/internal/domain"
	"drive/pkg/utils"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) domain.UserRepo {
	return &userRepo{db: db}
}

// 用户注册
func (r *userRepo) Register(user *domain.User) error {
	// 加密密码
	hashedPassword, err := utils.HashPassword(user.PassWord)
	if err != nil {
		return err
	}
	user.PassWord = hashedPassword

	return r.db.Create(user).Error
}

// 用户登录
func (r *userRepo) Logon(user *domain.User) error {
	// 先根据用户名查询用户
	var existingUser domain.User
	if err := r.db.Where("user_name = ?", user.UserName).First(&existingUser).Error; err != nil {
		return err
	}

	// 验证密码
	if !utils.CheckPasswordHash(user.PassWord, existingUser.PassWord) {
		return gorm.ErrRecordNotFound
	}

	// 将查询到的用户信息赋值给传入的user
	*user = existingUser
	return nil
}

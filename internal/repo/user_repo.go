package repo

import (
	"drive/internal/domain"

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
	return r.db.Create(user).Error
}

// 用户登录
func (r *userRepo) Logon(user *domain.User) error {
	return r.db.Where("user_name = ? AND pass_word = ?", user.UserName, user.PassWord).First(user).Error
}

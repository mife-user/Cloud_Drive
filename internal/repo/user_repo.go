package repo

import "drive/internal/domain"

type userRepo struct{}

func NewUserRepo() domain.UserRepo {
	return &userRepo{}
}

//用户注册
func (r *userRepo) Register(user *domain.User) error {

	return nil
}

//用户登录
func (r *userRepo) Logon(user *domain.User) error {

	return nil
}

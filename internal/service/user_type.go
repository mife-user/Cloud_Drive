package service

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/conf"
)

// UserServicer 用户处理器
type userServicer struct {
	userRepo domain.UserRepo
	config   *conf.Config
}

// RemixUser implements [domain.UserServicer].
func (s *userServicer) RemixUser(ctx context.Context, user *domain.User) error {
	panic("unimplemented")
}

// NewUserServicer 创建用户服务
func NewUserServicer(userRepo domain.UserRepo, config *conf.Config) domain.UserServicer {
	return &userServicer{
		userRepo: userRepo,
		config:   config,
	}
}

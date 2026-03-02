package service

import (
	"drive/internal/domain"
	"drive/pkg/conf"
)

// UserServicer 用户处理器
type userServicer struct {
	userRepo domain.UserRepo
	config   *conf.Config
}

// NewUserServicer 创建用户服务
func NewUserServicer(userRepo domain.UserRepo, config *conf.Config) domain.UserServicer {
	return &userServicer{
		userRepo: userRepo,
		config:   config,
	}
}

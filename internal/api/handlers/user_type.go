package handlers

import (
	"drive/internal/domain"
	"drive/pkg/conf"
)

// UserHandler 用户处理器
type UserHandler struct {
	userRepo domain.UserRepo
	config   *conf.Config
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userRepo domain.UserRepo, config *conf.Config) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		config:   config,
	}
}

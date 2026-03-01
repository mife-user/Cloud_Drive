package handlers

import (
	"drive/internal/domain"
	"drive/pkg/conf"
)

// UserHandler 用户处理器
type UserHandler struct {
	userServicer domain.UserServicer
	config       *conf.Config
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userServicer domain.UserServicer, config *conf.Config) *UserHandler {
	return &UserHandler{
		userServicer: userServicer,
		config:       config,
	}
}

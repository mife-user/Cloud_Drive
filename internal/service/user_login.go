package service

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/auth"
)

func (s *userServicer) Login(ctx context.Context, user *domain.User) (string, error) {
	var err error
	err = user.IsNullValue()
	if err != nil {
		return "", err
	}
	if err = s.userRepo.Logon(ctx, user); err != nil {
		return "", err
	}
	// 生成JWT token
	token, err := auth.GenerateToken(user.ID, user.Role, user.UserName, s.config.JWT.Secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

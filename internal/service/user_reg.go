package service

import (
	"context"
	"drive/internal/domain"
)

func (s *userServicer) Register(ctx context.Context, req *domain.User) error {
	var err error
	err = req.IsNullValue()
	if err != nil {
		return err
	}
	if err = s.userRepo.Register(ctx, req); err != nil {
		return err
	}
	return nil
}

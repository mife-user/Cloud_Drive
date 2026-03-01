package service

import (
	"context"
	"drive/pkg/errorer"
)

func (s *userServicer) GetUserHeadPath(ctx context.Context, username string) (string, error) {
	var err error
	var HeaderPath string
	if username == "" {
		err = errorer.New(errorer.ErrUserNameNotFound)
		return "", err
	}
	if HeaderPath, err = s.userRepo.GetUserHeadPath(ctx, username); err != nil {
		return "", err
	}
	return HeaderPath, nil
}

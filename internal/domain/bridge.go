package domain

import "context"

type UserRepo interface {
	Register(ctx context.Context, user *User) error
	Logon(user *User) error
}

type FileRepo interface {
	UploadFile(ctx context.Context, file *File) error
	ViewFile(ctx context.Context, userID string) ([]File, error)
}

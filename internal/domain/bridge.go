package domain

import "context"

type UserRepo interface {
	Register(ctx context.Context, user *User) error
	Logon(ctx context.Context, user *User) error
	RemixUser(ctx context.Context, user *User) error
}

type FileRepo interface {
	UploadFile(ctx context.Context, files []*File) error
	ViewFile(ctx context.Context, userID any) ([]File, error)
}

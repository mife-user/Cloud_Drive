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
	ShareFile(ctx context.Context, fileID uint, userID uint, owner string) (shareID string, accessKey string, err error)
	AccessShare(ctx context.Context, shareID string, accessKey string) (*File, error)
	UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error
	AddFavorite(ctx context.Context, userID uint, fileID uint, accessKey string) error
	RemoveFavorite(ctx context.Context, userID uint, fileID uint) error
	GetFavorites(ctx context.Context, userID uint) ([]File, error)
	CreateUploadTask(ctx context.Context, task *UploadTask) error
	GetUploadTaskByMD5(ctx context.Context, userID uint, fileMD5 string) (*UploadTask, error)
	GetUploadTaskByID(ctx context.Context, taskID uint) (*UploadTask, error)
	UpdateUploadTask(ctx context.Context, task *UploadTask) error
}

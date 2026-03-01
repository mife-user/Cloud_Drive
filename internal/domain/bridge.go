package domain

import (
	"context"
	"mime/multipart"
)

type UserRepo interface {
	Register(ctx context.Context, user *User) error
	Logon(ctx context.Context, user *User) error
	RemixUser(ctx context.Context, user *User) error
	UpdateHeader(ctx context.Context, header *UserHeader) error
	GetUserHeadPath(ctx context.Context, username string) (headPath string, err error)
}

type FileRepo interface {
	CheckUserSize(ctx context.Context, userID uint, totalSize int64) (int64, bool)
	UploadFile(ctx context.Context, files []*File, nowSize int64) error
	DeleteFile(ctx context.Context, userID uint, fileID uint) error
	DeleteFileForever(ctx context.Context, userID uint, fileID uint) error
	GetDeletedFiles(ctx context.Context, userID uint) ([]File, error)
	ViewFilesNote(ctx context.Context, userID uint) ([]File, error)
	ViewFile(ctx context.Context, fileID uint, userID uint) (*File, error)
	ShareFile(ctx context.Context, fileID uint, userID uint, owner string) (shareID string, accessKey string, err error)
	AccessShare(ctx context.Context, shareID string, accessKey string) (*File, error)
	UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error
	AddFavorite(ctx context.Context, userID uint, fileID uint) error
	RemoveFavorite(ctx context.Context, userID uint, fileID uint) error
	GetFavorites(ctx context.Context, userID uint) ([]File, error)
}
type UserServicer interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, user *User) (string, error)
	RemixUser(ctx context.Context, user *User) error
	UpdateHeader(ctx context.Context, header *UserHeader, fileHandler *multipart.FileHeader) error
	GetUserHeadPath(ctx context.Context, username string) (string, error)
}
type FileServicer interface {
	CheckUserSize(ctx context.Context, userID uint, totalSize int64) (int64, bool)
	UploadFile(ctx context.Context, files []*File, nowSize int64) error
	DeleteFile(ctx context.Context, userID uint, fileID uint) error
	DeleteFileForever(ctx context.Context, userID uint, fileID uint) error
	GetDeletedFiles(ctx context.Context, userID uint) ([]File, error)
	ViewFilesNote(ctx context.Context, userID uint) ([]File, error)
	ViewFile(ctx context.Context, fileID uint, userID uint) (*File, error)
	ShareFile(ctx context.Context, fileID uint, userID uint, owner string) (shareID string, accessKey string, err error)
	AccessShare(ctx context.Context, shareID string, accessKey string) (*File, error)
	UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error
	AddFavorite(ctx context.Context, userID uint, fileID uint) error
	RemoveFavorite(ctx context.Context, userID uint, fileID uint) error
	GetFavorites(ctx context.Context, userID uint) ([]File, error)
}

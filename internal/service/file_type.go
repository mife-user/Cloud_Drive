package service

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/conf"
)

type fileServicer struct {
	fileTypeRepo domain.FileRepo
	config       *conf.Config
}

// AccessShare implements [domain.FileServicer].
func (f *fileServicer) AccessShare(ctx context.Context, shareID string, accessKey string) (*domain.File, error) {
	panic("unimplemented")
}

// AddFavorite implements [domain.FileServicer].
func (f *fileServicer) AddFavorite(ctx context.Context, userID uint, fileID uint) error {
	panic("unimplemented")
}

// CheckUserSize implements [domain.FileServicer].
func (f *fileServicer) CheckUserSize(ctx context.Context, userID uint, totalSize int64) (int64, bool) {
	panic("unimplemented")
}

// DeleteFile implements [domain.FileServicer].
func (f *fileServicer) DeleteFile(ctx context.Context, userID uint, fileID uint) error {
	panic("unimplemented")
}

// DeleteFileForever implements [domain.FileServicer].
func (f *fileServicer) DeleteFileForever(ctx context.Context, userID uint, fileID uint) error {
	panic("unimplemented")
}

// GetDeletedFiles implements [domain.FileServicer].
func (f *fileServicer) GetDeletedFiles(ctx context.Context, userID uint) ([]domain.File, error) {
	panic("unimplemented")
}

// GetFavorites implements [domain.FileServicer].
func (f *fileServicer) GetFavorites(ctx context.Context, userID uint) ([]domain.File, error) {
	panic("unimplemented")
}

// RemoveFavorite implements [domain.FileServicer].
func (f *fileServicer) RemoveFavorite(ctx context.Context, userID uint, fileID uint) error {
	panic("unimplemented")
}

// ShareFile implements [domain.FileServicer].
func (f *fileServicer) ShareFile(ctx context.Context, fileID uint, userID uint, owner string) (shareID string, accessKey string, err error) {
	panic("unimplemented")
}

// UpdateFilePermissions implements [domain.FileServicer].
func (f *fileServicer) UpdateFilePermissions(ctx context.Context, fileID uint, userID uint, permissions string) error {
	panic("unimplemented")
}

// UploadFile implements [domain.FileServicer].
func (f *fileServicer) UploadFile(ctx context.Context, files []*domain.File, nowSize int64) error {
	panic("unimplemented")
}

// ViewFile implements [domain.FileServicer].
func (f *fileServicer) ViewFile(ctx context.Context, fileID uint, userID uint) (*domain.File, error) {
	panic("unimplemented")
}

// ViewFilesNote implements [domain.FileServicer].
func (f *fileServicer) ViewFilesNote(ctx context.Context, userID uint) ([]domain.File, error) {
	panic("unimplemented")
}

// NewFileTypeServicer 创建文件类型服务
func NewFileTypeServicer(fileTypeRepo domain.FileRepo, config *conf.Config) domain.FileServicer {
	return &fileServicer{
		fileTypeRepo: fileTypeRepo,
		config:       config,
	}
}

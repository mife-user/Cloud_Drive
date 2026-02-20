package repo

import (
	"context"
	"drive/internal/domain"
)

func (r *fileRepo) ViewFile(ctx context.Context, fileID uint) (*domain.File, error) {
	var file domain.File
	if err := r.db.Where("id = ?", fileID).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

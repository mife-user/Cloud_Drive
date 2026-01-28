package repo

import "drive/internal/domain"

type fileRepo struct{}

func NewFileRepo() domain.FileRepo {
	return &fileRepo{}
}

// 文件上传
func (r *fileRepo) UploadFile(file *domain.File) error {

	return nil
}


package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"
)

func (r *fileRepo) CreateUploadTask(ctx context.Context, task *domain.UploadTask) error {
	if err := r.db.Create(task).Error; err != nil {
		logger.Error("创建上传任务失败", logger.C(err))
		return err
	}
	logger.Info("创建上传任务成功")
	return nil
}

func (r *fileRepo) GetUploadTaskByMD5(ctx context.Context, userID uint, fileMD5 string) (*domain.UploadTask, error) {
	var task domain.UploadTask
	err := r.db.Where("user_id = ? AND file_md5 = ? AND status = ?", userID, fileMD5, 0).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *fileRepo) GetUploadTaskByID(ctx context.Context, taskID uint) (*domain.UploadTask, error) {
	var task domain.UploadTask
	err := r.db.Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *fileRepo) UpdateUploadTask(ctx context.Context, task *domain.UploadTask) error {
	if err := r.db.Save(task).Error; err != nil {
		logger.Error("更新上传任务失败", logger.C(err))
		return err
	}
	logger.Info("更新上传任务成功")
	return nil
}

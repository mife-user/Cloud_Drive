package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"fmt"
	"sync"
	"time"
)

// 文件上传
func (r *fileRepo) UploadFile(ctx context.Context, files []*domain.File, nowSize int64) error {
	var err error
	// 检查文件切片是否为空
	if len(files) == 0 {
		logger.Error(fmt.Sprintf("上传文件失败: %s", errorer.ErrEmptySlice))
		return fmt.Errorf(errorer.ErrEmptySlice)
	}
	// 上传文件到数据库
	if err = r.db.CreateInBatches(files, len(files)).Error; err != nil {
		logger.Error("上传文件失败", logger.C(err))
		return err
	}

	// 缓存文件信息
	userID := files[0].UserID
	userKey := fmt.Sprintf("files:%d", userID)
	workerPool := pool.NewPool(4) // 可根据系统配置调整大小
	workerPool.Start()
	var wg sync.WaitGroup
	for _, file := range files {
		f := file
		wg.Add(1)
		// 提交任务到协程池
		workerPool.Submit(func() {
			defer wg.Done()
			// 将文件信息序列化为 JSON
			fileJSON, err := exc.ExcFileToJSON(*f)
			if err != nil {
				logger.Error("序列化文件信息失败", logger.C(err))
				return // 继续处理其他文件，不影响整体操作
			}
			fileKey := fmt.Sprintf("file:%d", f.ID)
			// 缓存单个文件信息
			if err := r.rd.HSet(ctx, userKey, fileKey, fileJSON).Err(); err != nil {
				logger.Error("缓存文件信息失败", logger.C(err))
				return
			}
		})
	}
	// 等待所有任务完成后关闭协程池
	go func() {
		wg.Wait()
		workerPool.Stop()
	}()
	wg.Wait()
	r.rd.Expire(ctx, userKey, 3*time.Hour) // 设置缓存过期时间为3小时
	if err = r.db.Model(&domain.User{}).Where("id = ?", userID).Update("now_size", nowSize).Error; err != nil {
		logger.Error("更新用户空间失败", logger.C(err))
		return err
	}
	logger.Info("上传文件成功", logger.U("user_id", userID))
	return nil
}

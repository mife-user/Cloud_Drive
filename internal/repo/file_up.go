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
)

// 文件上传
func (r *fileRepo) UploadFile(ctx context.Context, files []*domain.File) error {
	// 检查文件切片是否为空
	if len(files) == 0 {
		logger.Error("上传文件失败: " + errorer.ErrEmptySlice)
		return fmt.Errorf(errorer.ErrEmptySlice)
	}
	// 上传文件到数据库
	if err := r.db.Create(files).Error; err != nil {
		logger.Error("上传文件失败", logger.C(err))
		return err
	}
	// 缓存文件信息
	userID := files[0].UserID
	userKey := fmt.Sprintf("files:%d", userID)
	pool := pool.NewPool(4) // 可根据系统配置调整大小
	pool.Start()
	var wg sync.WaitGroup
	for _, file := range files {
		f := file
		wg.Add(1)
		// 提交任务到协程池
		pool.Submit(func() {
			defer wg.Done()
			// 将文件信息序列化为 JSON
			fileJSON, err := exc.ExcFileToJSON(*f)
			if err != nil {
				logger.Error("序列化文件信息失败", logger.C(err))
				return // 继续处理其他文件，不影响整体操作
			}
			// 缓存单个文件信息
			if err := r.rd.HSet(ctx, userKey, f.FileName, fileJSON).Err(); err != nil {
				logger.Error("缓存文件信息失败", logger.C(err))
				return
			}
		})
	}
	// 等待所有任务完成
	wg.Wait()
	pool.Stop()
	logger.Info("上传文件成功")
	return nil
}

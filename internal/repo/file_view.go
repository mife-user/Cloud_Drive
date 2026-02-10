package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"encoding/json"
	"fmt"
	"sync"
)

// 查看文件
func (r *fileRepo) ViewFile(ctx context.Context, userID string) ([]domain.File, error) {
	var files []domain.File
	// 从缓存中查询文件信息
	userKey := fmt.Sprintf("files:%s", userID)            // 缓存键名
	fileJSONs, err := r.rd.HGetAll(ctx, userKey).Result() // 查询缓存中的所有文件信息
	if err != nil {
		logger.Error("查询缓存文件失败", logger.S("user_id", userID), logger.C(err))
		return nil, err
	}
	// 并发解析缓存中的文件信息
	pool := pool.NewPool(4)
	pool.Start()
	var wg sync.WaitGroup
	fileCh := make(chan domain.File, len(fileJSONs))
	errCh := make(chan error, 1)
	// 解析缓存中的文件信息
	for _, fileJSON := range fileJSONs {
		pool.Submit(func() {
			defer wg.Done()
			wg.Add(1)
			var file domain.File
			if err := json.Unmarshal([]byte(fileJSON), &file); err != nil {
				logger.Debug("解析缓存文件信息失败", logger.C(err))
				errCh <- err
				return
			}
			fileCh <- file
		})
	}
	// 等待所有文件解析完成
	wg.Wait()
	close(fileCh)
	close(errCh)
	// 检查是否有缓存错误
	if err, ok := <-errCh; ok {
		// 从数据库查询文件信息，仅返回文件名、大小、路径、权限和所有者
		if err = r.db.
			Select("file_name", "size", "path", "permissions", "owner").
			Where("user_id = ?", userID).
			Find(&files).Error; err != nil {
			logger.Error("查询文件失败", logger.S("user_id", userID), logger.C(err))
			return nil, err
		}
		logger.Info("从数据库查询文件成功")
		return files, nil
	}
	// 从通道中收集文件记录
	for file := range fileCh {
		files = append(files, file)
	}
	logger.Info("查询文件成功")
	return files, nil
}

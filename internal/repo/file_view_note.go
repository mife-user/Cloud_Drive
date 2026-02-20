package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"fmt"
	"sync"
)

// 查看文件
func (r *fileRepo) ViewFilesNote(ctx context.Context, userID uint) ([]domain.File, error) {
	var err error
	var result bool = false
	var files []domain.File
	var filesNew []domain.File
	userIDStr := fmt.Sprintf("%d", userID)
	// 从缓存中查询文件信息
	userKey := fmt.Sprintf("files:%s", userIDStr)         // 缓存键名
	fileJSONs, err := r.rd.HGetAll(ctx, userKey).Result() // 查询缓存中的所有文件信息
	if err != nil {
		logger.Error("查询缓存文件失败", logger.S("user_id", userIDStr), logger.C(err))
		return nil, err
	}
	// 并发解析缓存中的文件信息
	workerPool := pool.NewPool(4)
	workerPool.Start()
	var wg sync.WaitGroup
	fileCh := make(chan domain.File, len(fileJSONs))
	// 解析缓存中的文件信息
	for _, fileJSON := range fileJSONs {
		wg.Add(1)
		workerPool.Submit(func() {
			defer wg.Done()
			var file domain.File
			if err = exc.ExcJSONToFile(fileJSON, &file); err != nil {
				logger.Debug("解析缓存文件信息失败", logger.C(err))
				fileCh <- file
			}
			fileCh <- file
		})
	}
	// 等待所有任务完成
	workerPool.Stop()
	wg.Wait()
	close(fileCh)
	// 从通道中收集文件记录
	for file := range fileCh {
		if file.FileName == "" {
			result = true
			break
		}
		files = append(files, file)
	}
	if result {
		// 从数据库查询文件信息，仅返回文件名、大小、路径、权限和所有者
		if err = r.db.
			Select("file_name", "size", "path", "permissions", "owner").
			Where("user_id = ?", userID).
			Find(&filesNew).Error; err != nil {
			logger.Error("查询文件失败", logger.S("user_id", userIDStr), logger.C(err))
			return nil, err
		}
		logger.Info("从数据库查询文件成功")
	}
	logger.Info("查询文件成功")
	return files, nil
}

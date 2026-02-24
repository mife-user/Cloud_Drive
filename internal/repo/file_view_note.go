package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/cache"
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
	// 从缓存中查询文件信息
	userKey := fmt.Sprintf("files:%d", userID)
	fileJSONs, err := r.rd.HGetAll(ctx, userKey).Result()
	if err != nil {
		logger.Error("查询缓存文件信息失败", logger.C(err))
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
		fileJSONCopy := fileJSON
		workerPool.Submit(func() {
			defer wg.Done()
			var file domain.File
			if err = exc.ExcJSONToFile(fileJSONCopy, &file); err != nil {
				logger.Debug("解析缓存文件信息失败", logger.C(err))
				return
			}
			if file.DeletedAt.Valid {
				return
			}
			fileCh <- file
		})
	}
	// 等待所有任务完成后关闭通道和协程池
	go func() {
		wg.Wait()
		close(fileCh)
		workerPool.Stop()
	}()
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
			Select("file_name", "size", "path", "permissions", "owner").Where("user_id = ?", userID).
			Find(&filesNew).Error; err != nil {
			logger.Error("查询文件失败", logger.U("user_id", userID), logger.C(err))
			return nil, err
		}
		// 更新缓存中的文件信息
		for _, file := range filesNew {
			fileJSON, err := exc.ExcFileToJSON(file)
			if err != nil {
				logger.Error("序列化文件信息失败", logger.C(err))
				continue
			}
			if err = r.rd.HSet(ctx, userKey, fmt.Sprintf("file:%d", file.ID), fileJSON).Err(); err != nil {
				logger.Error("缓存文件信息失败", logger.C(err))
				continue
			}
			ttl := cache.FileCacheConfig.RandomTTL()
			if err = r.rd.Expire(ctx, userKey, ttl).Err(); err != nil {
				logger.Error("设置缓存过期时间失败", logger.C(err))
				continue
			}
		}
		logger.Info("从数据库查询文件成功", logger.U("user_id", userID))
		return filesNew, nil
	}
	logger.Info("查询文件成功", logger.U("user_id", userID))
	return files, nil
}

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

func (r *fileRepo) GetDeletedFiles(ctx context.Context, userID uint) ([]domain.File, error) {
	var err error
	var files []domain.File
	//缓存查询用户删除的文件
	userKey := fmt.Sprintf("files:%d", userID)
	mapcmd := r.rd.HGetAll(ctx, userKey)
	if err = mapcmd.Err(); err != nil {
		// 查询用户删除的文件
		if err = r.db.Unscoped().Where("deleted_at IS NOT NULL AND user_id = ?", userID).Find(&files).Error; err != nil {
			logger.Error("查询用户删除的文件失败", logger.U("user_id", userID), logger.C(err))
			return nil, err
		}
	} else {
		fileJSONs := mapcmd.Val()
		var wg sync.WaitGroup
		fileChan := make(chan domain.File, len(fileJSONs))
		workerPool := pool.NewPool(4)
		workerPool.Start()
		for _, fileJSON := range fileJSONs {
			f := fileJSON
			wg.Add(1)
			workerPool.Submit(func() {
				defer wg.Done()
				var file domain.File
				if err = exc.ExcJSONToFile(f, &file); err != nil {
					logger.Error("反序列化文件失败", logger.C(err))
					return
				}
				if file.DeletedAt.Valid {
					fileChan <- file
				} else {
					return
				}
			})
		}
		go func() {
			wg.Wait()
			close(fileChan)
			workerPool.Stop()
		}()
		for file := range fileChan {
			files = append(files, file)
		}
	}
	logger.Info("查询用户删除的文件成功", logger.U("user_id", userID))
	return files, nil
}

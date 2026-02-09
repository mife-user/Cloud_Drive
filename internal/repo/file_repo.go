package repo

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type fileRepo struct {
	db *gorm.DB
	rd *redis.Client
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB, rd *redis.Client) domain.FileRepo {
	return &fileRepo{db: db, rd: rd}
}

// 文件上传
func (r *fileRepo) UploadFile(ctx context.Context, files []*domain.File) error {
	// 检查文件切片是否为空
	if len(files) == 0 {
		logger.Error("上传文件失败: " + errorer.ErrEmptySlice)
		return fmt.Errorf(errorer.ErrEmptySlice)
	}
	// 上传文件到数据库
	if err := r.db.Create(files).Error; err != nil {
		logger.Error("上传文件失败", zap.Error(err))
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
			fileJSON, err := json.Marshal(f)
			if err != nil {
				logger.Error("序列化文件信息失败", zap.Error(err))
				return // 继续处理其他文件，不影响整体操作
			}
			// 缓存单个文件信息
			if err := r.rd.HSet(ctx, userKey, f.FileName, string(fileJSON)).Err(); err != nil {
				logger.Error("缓存文件信息失败", zap.Error(err))
				return
			}
		})
	}
	// 等待所有任务完成
	wg.Wait()
	pool.Stop()
	logger.Debug("上传文件成功")
	return nil
}

// 查看文件
func (r *fileRepo) ViewFile(ctx context.Context, userID string) ([]domain.File, error) {
	var files []domain.File
	// 从缓存中查询文件信息
	userKey := fmt.Sprintf("files:%s", userID)
	fileJSONs, err := r.rd.HGetAll(ctx, userKey).Result()
	if err != nil {
		logger.Error("查询缓存文件失败", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	pool := pool.NewPool(4) // 可根据系统配置调整大小
	pool.Start()
	var wg sync.WaitGroup
	fileCh := make(chan domain.File, len(fileJSONs))
	// 解析缓存中的文件信息
	for _, fileJSON := range fileJSONs {
		pool.Submit(func() {
			defer wg.Done()
			wg.Add(1)
			var file domain.File
			if err := json.Unmarshal([]byte(fileJSON), &file); err != nil {
				logger.Error("解析缓存文件信息失败", zap.Error(err))
			}
			fileCh <- file
		})
	}
	go func() {
		wg.Wait()
		close(fileCh)
	}()
	// 从通道中收集文件记录
	for file := range fileCh {
		files = append(files, file)
	}
	// 从数据库查询文件信息，仅返回文件名、大小、路径、权限和所有者
	if err := r.db.Select("file_name", "size", "path", "permissions", "owner").
		Where("user_id = ?", userID).
		Find(&files).Error; err != nil {
		logger.Error("查询文件失败", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	logger.Debug("查询文件成功")
	return files, nil
}

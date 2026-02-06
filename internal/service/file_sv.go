package service

import (
	"drive/internal/domain"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"drive/pkg/utils"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SaveFiles(c *gin.Context, files []*multipart.FileHeader, fileRecords *[]*domain.File, userID uint) error {
	// 保存文件
	pool := pool.NewPool(4)
	pool.Start()
	for _, header := range files {
		// 提交任务到线程池
		pool.Submit(func() {
			fileRecord, err := utils.SaveFile(header, userID)
			if err != nil {
				logger.Error("保存文件失败: %v", zap.Error(err))
				return
			}
			// 保存文件记录到切片
			*fileRecords = append(*fileRecords, fileRecord)
		})
	}
	pool.Stop()
	return nil
}

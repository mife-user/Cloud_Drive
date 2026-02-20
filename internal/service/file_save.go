package service

import (
	"drive/internal/domain"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"drive/pkg/utils"
	"mime/multipart"
	"sync"
)

// SaveFiles 保存文件
func SaveFiles(files []*multipart.FileHeader, userID uint, userName string, userRole string, filekey *domain.File) (*[]*domain.File, error) {
	if userRole != "VIP" {
		filekey.Size = 1073741824 // 普通用户文件大小限制为1GB
	} else {
		filekey.Size = 2147483648 // 会员用户文件大小限制为2GB
	}
	//初始化文件
	filekey.Owner = userName
	filekey.UserID = userID
	// 创建文件记录通道
	recordCh := make(chan *domain.File, len(files))
	// 保存文件
	pool := pool.NewPool(4) // 可根据系统配置调整大小
	pool.Start()
	var wg sync.WaitGroup
	for _, header := range files {
		h := header
		wg.Add(1)
		// 提交任务到协程池
		pool.Submit(func() {
			defer wg.Done()
			fileRecord, err := utils.SaveFile(h, filekey)
			if err != nil {
				logger.Error("保存文件失败: %v", logger.C(err))
				return
			}
			// 将文件记录发送到通道
			recordCh <- fileRecord
		})
	}
	wg.Wait()
	close(recordCh)
	// 从通道中收集文件记录
	fileRecords := make([]*domain.File, 0, len(files))
	for record := range recordCh {
		fileRecords = append(fileRecords, record)
	}

	// 关闭协程池
	pool.Stop()
	return &fileRecords, nil
}

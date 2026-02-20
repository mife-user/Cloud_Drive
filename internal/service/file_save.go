package service

import (
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"drive/pkg/utils"
	"mime/multipart"
	"sync"
)

// SaveFiles 保存文件
func SaveFiles(files []*multipart.FileHeader, userID any, userName any, userRole any, filekey *domain.File) (*[]*domain.File, error) {
	// 检查用户角色是否为会员
	userRoleStr, ok := userRole.(string)
	if !ok {
		logger.Error("userRole类型转换失败")
		return nil, errorer.New(errorer.ErrTypeError)
	}
	if userRoleStr != "VIP" {
		filekey.Size = 1024 * 1024 * 1024 // 普通用户文件大小限制为1GB
	} else {
		filekey.Size = 1024 * 1024 * 2048 // 会员用户文件大小限制为2GB
	}
	// 转换userID为uint类型
	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("userID类型转换失败")
		return nil, errorer.New(errorer.ErrTypeError)
	}
	// 转换userName为string类型
	userNameStr, ok := userName.(string)
	if !ok {
		logger.Error("userName类型转换失败")
		return nil, errorer.New(errorer.ErrTypeError)
	}
	//初始化文件
	filekey.Owner = userNameStr
	filekey.UserID = userIDUint
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

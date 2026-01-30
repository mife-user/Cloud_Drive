package utils

import (
	"drive/internal/domain"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveFile(header *multipart.FileHeader, file multipart.File, userID any) (*domain.File, error) {
	// 创建存储目录结构
	storageDir := fmt.Sprintf("./storage/%v/%s", userID, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(storageDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 保存文件
	if _, err := dst.ReadFrom(file); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 创建文件记录
	fileRecord := &domain.File{
		FileName:    header.Filename,
		Size:        header.Size,
		Path:        filePath,
		UserID:      userID.(uint),
		Permissions: "private",
	}
	return fileRecord, nil
}

package save

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SaveFile 保存文件
func SaveFile(header *multipart.FileHeader, size int64, userID uint) (string, int64, string, error) {
	// 检查文件大小是否超过限制
	if header.Size > size {
		return "", 0, "", fmt.Errorf("单个文件大小超过限制: %v", header.Size)
	}
	// 打开文件
	file, err := header.Open()
	if err != nil {
		return "", 0, "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 提取目录路径和文件名
	dirPath := filepath.Dir(header.Filename)   // 包含子目录路径
	fileName := filepath.Base(header.Filename) // 文件名
	//查看路径是否包含..
	if strings.Contains(dirPath, "..") || strings.Contains(fileName, "..") {
		return "", 0, "", fmt.Errorf("文件名包含无效路径: %s", header.Filename)
	}
	// 创建存储目录结构
	storageBase := fmt.Sprintf("./storage/%v", userID)
	storageDir := filepath.Join(storageBase, dirPath)
	// 检查目录是否存在，不存在则创建
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", 0, "", fmt.Errorf("创建存储目录失败: %w", err)
	}

	// 生成唯一文件名（处理命名冲突）
	ext := filepath.Ext(fileName)
	baseName := fileName[:len(fileName)-len(ext)]
	finalFileName := fileName
	tempPath := filepath.Join(storageDir, finalFileName)

	// 检查文件是否存在，存在则添加时间戳后缀
	if _, err := os.Stat(tempPath); err == nil {
		finalFileName = fmt.Sprintf("%s_%d%s", baseName, time.Now().UnixNano(), ext)
		tempPath = filepath.Join(storageDir, finalFileName)
	}

	// 创建目标文件
	dst, err := os.Create(tempPath)
	if err != nil {
		return "", 0, "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 保存文件
	if _, err := dst.ReadFrom(file); err != nil {
		return "", 0, "", fmt.Errorf("保存文件失败: %w", err)
	}

	return header.Filename, header.Size, tempPath, nil
}

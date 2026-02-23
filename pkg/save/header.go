package save

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveHeader 保存用户头像
func SaveHeader(ctx context.Context, fileHeader *multipart.FileHeader, userID uint, rating int) (uint, string, error) {
	var err error
	//获取后缀
	ext := filepath.Ext(fileHeader.Filename)
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return 0, "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	// 保存文件
	headerPath := fmt.Sprintf("./header/%v", userID)
	// 检查目录是否存在，不存在则创建
	if err = os.MkdirAll(headerPath, 0755); err != nil {
		return 0, "", fmt.Errorf("创建存储目录失败: %w", err)
	}
	// 生成唯一文件名（处理命名冲突）
	finalHeaderPath := fmt.Sprintf("%s_%d%s", headerPath, time.Now().UnixNano(), ext)
	// 创建目标文件
	dst, err := os.Create(finalHeaderPath)
	if err != nil {
		return 0, "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()
	//限制器
	limiter := newReadLimit(ctx, file, rating, 1024*1024)
	// 保存文件
	if _, err = dst.ReadFrom(limiter); err != nil {
		return 0, "", fmt.Errorf("保存文件失败: %w", err)
	}
	return userID, headerPath, nil
}

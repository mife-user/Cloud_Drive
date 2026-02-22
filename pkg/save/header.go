package save

import (
	"fmt"
	"mime/multipart"
	"os"
)

// SaveHeader 保存用户头像
func SaveHeader(fileHeader *multipart.FileHeader, userID uint) (uint, string, error) {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return 0, "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	// 保存文件
	headerPath := fmt.Sprintf("./header/%v", userID)
	// 检查目录是否存在，不存在则创建
	if err := os.MkdirAll(headerPath, 0755); err != nil {
		return 0, "", fmt.Errorf("创建存储目录失败: %w", err)
	}
	// 创建目标文件
	dst, err := os.Create(headerPath)
	if err != nil {
		return 0, "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()
	// 保存文件
	if _, err := dst.ReadFrom(file); err != nil {
		return 0, "", fmt.Errorf("保存文件失败: %w", err)
	}
	return userID, headerPath, nil
}

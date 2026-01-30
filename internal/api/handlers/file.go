package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"drive/internal/domain"
	"drive/pkg/conf"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileRepo domain.FileRepo
	config   *conf.Config
}

// NewFileHandler 创建文件处理器
func NewFileHandler(fileRepo domain.FileRepo, config *conf.Config) *FileHandler {
	return &FileHandler{
		fileRepo: fileRepo,
		config:   config,
	}
}

// UploadFile 文件上传
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
		return
	}
	defer file.Close()

	// 从认证中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	// 创建存储目录结构
	storageDir := fmt.Sprintf("./storage/%v/%s", userID, time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建存储目录失败: " + err.Error()})
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(storageDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败: " + err.Error()})
		return
	}
	defer dst.Close()

	// 保存文件
	if _, err := dst.ReadFrom(file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	// 创建文件记录
	fileRecord := &domain.File{
		FileName:    header.Filename,
		Size:        header.Size,
		Path:        filePath,
		UserID:      userID.(uint),
		Permissions: "private",
	}

	// 保存文件记录到数据库
	if err := h.fileRepo.UploadFile(fileRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"file": gin.H{
			"file_name": fileRecord.FileName,
			"size":      fileRecord.Size,
			"path":      fileRecord.Path,
			"user_id":   fileRecord.UserID,
		},
	})
}

package handlers

import (
	"net/http"

	"drive/internal/domain"
	"drive/pkg/conf"
	"drive/pkg/utils"

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
	// 保存文件
	fileRecord, err := utils.SaveFile(header, file, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
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

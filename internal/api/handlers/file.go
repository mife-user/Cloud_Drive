package handlers

import (
	"net/http"

	"drive/internal/domain"
	"drive/pkg/conf"
	"drive/pkg/pool"
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
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
		return
	}
	// 获取上传的文件头
	files := file.File["files"]

	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	// 保存文件
	var fileRecords []*domain.File
	pool := pool.NewPool(4)
	pool.Start()
	for _, header := range files {
		// 提交任务到线程池
		pool.Submit(func() {
			fileRecord, err := utils.SaveFile(header, userID.(uint))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
				return
			}
			// 保存文件记录到切片
			fileRecords = append(fileRecords, fileRecord)
		})
	}
	pool.Stop()

	// 保存文件记录到数据库
	if err := h.fileRepo.UploadFile(c, fileRecords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
	})
}

// 查看所有文件
func (h *FileHandler) ViewFiles(c *gin.Context) {

	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	// 查看文件
	files, err := h.fileRepo.ViewFile(c, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查看文件失败: " + err.Error()})
		return
	}
	// 返回文件列表
	c.JSON(http.StatusOK, gin.H{
		"message": "查看成功",
		"files":   files,
	})
}

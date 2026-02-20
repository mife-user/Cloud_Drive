package handlers

import (
	"drive/internal/api/dtos"
	"drive/internal/service"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadFile 文件上传
func (h *FileHandler) UploadFile(c *gin.Context) {
	logger.Info("开始处理文件上传请求")
	defer logger.Info("文件上传请求处理完成")
	// 绑定 JSON 请求体到 FileDtos
	var fileDto dtos.FileDtos
	// 获取上传的文件
	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
		return
	}
	fileDto.Permissions = file.Value["permissions"][0]
	// 获取上传的文件头
	files := file.File["files"]
	// 获取当前登录用户ID
	userID, existsID := c.Get("user_id")
	if !existsID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	// 获取当前登录用户名
	userName, existsNM := c.Get("user_name")
	if !existsNM {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	userRole, existsRole := c.Get("role")
	if !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	// 保存文件
	fileRecords, err := service.SaveFiles(files, userID, userName, userRole, fileDto.ToDMFile())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	// 保存文件记录到数据库
	if len(*fileRecords) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: 没有文件记录"})
		return
	}
	if err := h.fileRepo.UploadFile(c, *fileRecords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
	})
}

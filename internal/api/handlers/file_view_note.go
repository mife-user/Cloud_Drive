package handlers

import (
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 查看所有文件
func (h *FileHandler) ViewFilesNote(c *gin.Context) {
	logger.Info("开始处理查看所有文件请求")
	defer logger.Info("查看所有文件请求处理完成")

	// 获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户ID类型错误"})
		return
	}
	// 查看文件
	files, err := h.fileRepo.ViewFilesNote(c, userIDUint)
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

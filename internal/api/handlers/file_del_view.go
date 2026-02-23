package handlers

import (
	"context"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) ViewDeletedFiles(c *gin.Context) {
	logger.Info("开始处理查看用户删除文件请求")
	defer logger.Info("查看用户删除文件请求处理完成")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	//解析用户ID
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	//查询用户删除的文件
	files, err := h.fileRepo.GetDeletedFiles(ctx, userIDUint)
	if err != nil {
		logger.Error("查询用户删除的文件失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户删除的文件失败"})
		return
	}
	// 返回文件列表
	c.JSON(http.StatusOK, gin.H{
		"message": "查看成功",
		"files":   files,
	})
}

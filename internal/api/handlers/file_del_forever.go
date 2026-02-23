package handlers

import (
	"context"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 永久删除文件
func (h *FileHandler) DeleteFileForever(c *gin.Context) {
	logger.Info("开始处理永久删除文件请求")
	defer logger.Info("永久删除文件请求处理完成")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Error("用户ID不存在")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户ID不存在"})
		return
	}
	// 从URL参数获取文件ID
	fileID := c.Param("file_id")
	//解析文件ID
	fileIDUint, err := exc.StrToUint(fileID)
	if err != nil {
		logger.Error("文件ID格式错误", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID格式错误"})
		return
	}
	//解析用户id
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		logger.Error("用户ID格式错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID格式错误"})
		return
	}
	// 调用服务层删除文件
	if err := h.fileRepo.DeleteFileForever(ctx, userIDUint, fileIDUint); err != nil {
		logger.Error("删除文件失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文件永久删除成功"})
}

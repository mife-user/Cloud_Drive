package handlers

import (
	"drive/internal/service"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) UpdateHeader(c *gin.Context) {
	logger.Info("开始处理更新用户头像请求")
	defer logger.Info("更新用户头像请求处理完成")
	// 获取用户ID
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	// 获取用户上传的文件
	file, err := c.FormFile("header")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	// 检查userID是否为uint类型
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	// 调用服务层更新用户头像
	fileRecord, err := service.UpdateHeader(file, userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户头像失败"})
		switch err.Error() {
		case errorer.ErrFileSizeExceeded:
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过10MB"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "未知错误"})
		}
		return
	}
	if err := h.userRepo.UpdateHeader(c, fileRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户头像失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户头像更新成功"})
}

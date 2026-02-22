package handlers

import (
	"context"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RemoveFavorite 取消文件收藏
func (h *FileHandler) RemoveFavorite(c *gin.Context) {
	logger.Info("开始处理取消收藏请求")
	// 设置合理的超时时间，取消收藏操作涉及数据库和缓存
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	fileID, err := exc.StrToUint(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID格式错误"})
		return
	}

	userID, existsID := c.Get("user_id")
	if !existsID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	if err := h.fileRepo.RemoveFavorite(ctx, userIDUint, fileID); err != nil {

		switch err.Error() {
		case errorer.ErrFavoriteNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "收藏不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		}
		return
	}

	logger.Info("取消收藏成功", logger.S("file_id", fileIDStr))
	c.JSON(http.StatusOK, gin.H{
		"message": "取消收藏成功",
	})
}

package handlers

import (
	"context"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetFavorites 获取收藏列表
func (h *FileHandler) GetFavorites(c *gin.Context) {
	logger.Info("开始处理获取收藏列表请求")
	// 设置合理的超时时间，获取收藏列表涉及数据库查询
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, existsID := c.Get("user_id")
	if !existsID {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		logger.Error("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	files, err := h.fileRepo.GetFavorites(ctx, userIDUint)
	if err != nil {
		logger.Error("获取收藏列表失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败"})
		return
	}

	logger.Info("获取收藏列表成功", logger.S("count", fmt.Sprintf("%d", len(files))))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"files":   files,
	})
}

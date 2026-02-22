package handlers

import (
	"context"
	"drive/internal/api/dtos"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AddFavorite 添加文件收藏
func (h *FileHandler) AddFavorite(c *gin.Context) {
	logger.Info("开始处理添加收藏请求")
	// 设置合理的超时时间，收藏操作涉及数据库和缓存
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var req dtos.FavoriteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("绑定请求参数失败", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

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

	if err := h.fileRepo.AddFavorite(ctx, userIDUint, req.FileID); err != nil {
		logger.Error("添加收藏失败", logger.S("file_id", fmt.Sprintf("%d", req.FileID)), logger.C(err))
		switch err.Error() {
		case errorer.ErrFileNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		case errorer.ErrNotFileOwner:
			c.JSON(http.StatusForbidden, gin.H{"error": "非文件所有者，无法收藏"})
		case errorer.ErrFavoriteExist:
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件已收藏"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加收藏失败"})
		}
		return
	}

	logger.Info("添加收藏成功", logger.S("file_id", fmt.Sprintf("%d", req.FileID)))
	c.JSON(http.StatusOK, gin.H{
		"message": "收藏成功",
	})
}

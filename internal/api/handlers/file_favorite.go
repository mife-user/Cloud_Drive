package handlers

import (
	"drive/internal/api/dtos"
	"drive/internal/service"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddFavorite 添加文件收藏
func (h *FileHandler) AddFavorite(c *gin.Context) {
	logger.Info("开始处理添加收藏请求")
	defer logger.Info("添加收藏请求处理完成")

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

	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	if err := h.fileRepo.AddFavorite(c, userIDUint, req.FileID, req.AccessKey); err != nil {
		logger.Error("添加收藏失败", logger.S("file_id", fmt.Sprintf("%d", req.FileID)), logger.C(err))
		switch err.Error() {
		case errorer.ErrFileNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		case errorer.ErrNotFileOwner:
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限收藏文件，请确保您是文件所有者或持有有效的访问密钥"})
		case errorer.ErrFavoriteExist:
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件已收藏"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加收藏失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "收藏成功",
	})
}

// RemoveFavorite 取消文件收藏
func (h *FileHandler) RemoveFavorite(c *gin.Context) {
	logger.Info("开始处理取消收藏请求")
	defer logger.Info("取消收藏请求处理完成")

	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		logger.Error("文件ID为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	fileID, err := service.ParseID(fileIDStr)
	if err != nil {
		logger.Error("文件ID格式错误", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID格式错误"})
		return
	}

	userID, existsID := c.Get("user_id")
	if !existsID {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	if err := h.fileRepo.RemoveFavorite(c.Request.Context(), userIDUint, fileID); err != nil {
		logger.Error("取消收藏失败", logger.S("file_id", fileIDStr), logger.C(err))
		switch err.Error() {
		case errorer.ErrFavoriteNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "收藏不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "取消收藏成功",
	})
}

// GetFavorites 获取收藏列表
func (h *FileHandler) GetFavorites(c *gin.Context) {
	logger.Info("开始处理获取收藏列表请求")
	defer logger.Info("获取收藏列表请求处理完成")

	userID, existsID := c.Get("user_id")
	if !existsID {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	files, err := h.fileRepo.GetFavorites(c.Request.Context(), userIDUint)
	if err != nil {
		logger.Error("获取收藏列表失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"files":   files,
	})
}

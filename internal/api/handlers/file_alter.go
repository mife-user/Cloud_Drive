package handlers

import (
	"context"
	"drive/internal/api/dtos/request"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateFilePermissions 更新文件权限
func (h *FileHandler) UpdateFilePermissions(c *gin.Context) {
	logger.Info("开始处理更新文件权限请求")
	// 设置合理的超时时间，文件权限更新涉及数据库和缓存操作
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		logger.Error("文件ID为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	fileID, err := exc.StrToUint(fileIDStr)
	if err != nil {
		logger.Error("文件ID格式错误", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID格式错误"})
		return
	}

	var req request.FileDtos
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

	if err := h.fileRepo.UpdateFilePermissions(ctx, fileID, userIDUint, req.Permissions); err != nil {
		logger.Error("更新文件权限失败", logger.S("file_id", fileIDStr), logger.C(err))
		switch err.Error() {
		case errorer.ErrFileNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		case errorer.ErrNotFileOwner:
			c.JSON(http.StatusForbidden, gin.H{"error": "非文件所有者，无法更新权限"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		}
		return
	}

	logger.Info("更新文件权限成功", logger.S("file_id", fileIDStr))
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

package handlers

import (
	"context"
	"drive/internal/api/dtos/request"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) UpdateFilePermissions(c *gin.Context) {
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

	var req request.FilePermissionsDT
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
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

	if err := h.fileServicer.UpdateFilePermissions(ctx, fileID, userIDUint, req.Permissions); err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

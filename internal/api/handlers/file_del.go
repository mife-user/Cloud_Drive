package handlers

import (
	"context"
	"drive/pkg/exc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) DeleteFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID类型错误"})
		return
	}

	fileIDStr := c.Param("file_id")
	fileIDUint, err := exc.StrToUint(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.fileServicer.DeleteFile(ctx, userIDUint, fileIDUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功", "回收站保存时间": "24小时"})
}

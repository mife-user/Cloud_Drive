package handlers

import (
	"context"
	"drive/internal/api/dtos/request"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) AccessShare(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	shareID := c.Param("share_id")
	if shareID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分享ID不能为空"})
		return
	}

	var req request.AccessShareDT
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	file, err := h.fileServicer.AccessShare(ctx, shareID, req.AccessKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":   file.ID,
		"file_name": file.FileName,
		"size":      file.Size,
		"path":      file.Path,
		"owner":     file.Owner,
	})
}

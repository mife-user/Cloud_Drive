package handlers

import (
	"drive/internal/api/dtos"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccessShare 访问文件分享
func (h *FileHandler) AccessShare(c *gin.Context) {
	logger.Info("开始处理访问分享请求")

	shareID := c.Param("share_id")
	if shareID == "" {
		logger.Error("分享ID为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "分享ID不能为空"})
		return
	}

	var req dtos.AccessShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("绑定请求参数失败", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	file, err := h.fileRepo.AccessShare(c.Request.Context(), shareID, req.AccessKey)
	if err != nil {
		logger.Error("访问分享失败", logger.S("share_id", shareID), logger.C(err))
		switch err.Error() {
		case errorer.ErrShareNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "分享不存在"})
		case errorer.ErrInvalidAccessKey:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "访问密钥无效"})
		case errorer.ErrShareExpired:
			c.JSON(http.StatusGone, gin.H{"error": "分享已过期"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "访问失败"})
		}
		return
	}

	logger.Info("访问分享成功", logger.S("share_id", shareID), logger.S("file_path", file.Path))
	c.File(file.Path)
}

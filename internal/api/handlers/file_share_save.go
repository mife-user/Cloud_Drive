package handlers

import (
	"drive/internal/api/dtos"
	"drive/internal/service"
	"drive/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShareFile 分享文件
func (h *FileHandler) ShareFile(c *gin.Context) {
	logger.Info("开始处理分享文件请求")

	var req dtos.ShareFileRequest
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

	userName, existsNM := c.Get("user_name")
	if !existsNM {
		logger.Error("获取用户名失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	// 转换userID和userName为uint和string类型
	userIDUint, userNameStr, err := service.ExchangeFile(userID, userName)
	if err != nil {
		logger.Error("转换用户ID和用户名失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 分享文件
	shareID, accessKey, err := h.fileRepo.ShareFile(c, req.FileID, userIDUint, userNameStr)
	if err != nil {
		logger.Error("分享文件失败", logger.S("file_id", fmt.Sprintf("%d", req.FileID)), logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 构建分享URL
	shareURL := fmt.Sprintf("http://%s/api/file/share/%s", c.Request.Host, shareID)

	logger.Info("分享文件成功", logger.S("share_id", shareID), logger.S("file_id", fmt.Sprintf("%d", req.FileID)))
	c.JSON(http.StatusOK, gin.H{
		"message":    "分享成功",
		"share_url":  shareURL,
		"access_key": accessKey,
	})
}

package handlers

import (
	"context"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) DeleteFile(c *gin.Context) {
	logger.Info("开始处理文件删除请求")
	defer logger.Info("文件删除请求处理完成")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	//获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	//解析id
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		logger.Error("用户ID格式错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID类型错误"})
		return
	}
	fileIDStr := c.Param("file_id")
	fileIDUint, err := exc.StrToUint(fileIDStr)
	if err != nil {
		logger.Error("文件ID格式错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//删除文件
	if err := h.fileRepo.DeleteFile(ctx, userIDUint, fileIDUint); err != nil {
		logger.Error("删除文件失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功", "回收站保存时间": "1小时"})
}

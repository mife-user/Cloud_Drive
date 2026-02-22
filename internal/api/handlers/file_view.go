package handlers

import (
	"context"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 查看单个文件
func (h *FileHandler) ViewFile(c *gin.Context) {
	logger.Info("开始处理查看单个文件请求")
	defer logger.Info("查看单个文件请求处理完成")
	// 设置合理的超时时间，查看文件涉及数据库查询和文件读取
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	// 获取文件ID
	fileID := c.Param("file_id")
	//将文件ID转换为uint类型
	fileIDUint, err := exc.StrToUint(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID格式错误"})
		return
	}
	// 查看文件
	file, err := h.fileRepo.ViewFile(ctx, fileIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查看文件失败: " + err.Error()})
		return
	}
	// 返回文件信息
	c.File(file.Path)
}

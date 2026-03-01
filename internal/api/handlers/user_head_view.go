package handlers

import (
	"context"
	"drive/internal/api/dtos/response"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetHeader(c *gin.Context) {
	logger.Info("开始处理用户信息修改请求")
	defer logger.Info("用户信息修改请求处理完成")
	//设置过期时间
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	// 获取用户名
	username := c.Param("username")
	if username == "" {
		logger.Debug("获取用户信息失败" + errorer.ErrUserIDNotFound)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorer.ErrUserIDNotFound})
		return
	}
	// 查询用户头像
	headPath, err := h.userServicer.GetUserHeadPath(ctx, username)
	if err != nil {
		logger.Error("查询用户头像失败", logger.S("username", username), logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户头像失败"})
		return
	}
	headPathRSP := response.ToDTHeaderPath(headPath)
	// 返回用户头像
	c.JSON(http.StatusOK, gin.H{
		"message": "查询用户头像成功",
		"data":    headPathRSP,
	})
}

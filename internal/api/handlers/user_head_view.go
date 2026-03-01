package handlers

import (
	"context"
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
	// 获取用户名称
	userName := c.Param("user_name")
	if userName == "" {
		logger.Debug("获取用户信息失败" + errorer.ErrUserNameNotFound)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorer.ErrUserNameNotFound})
		return
	}
	// 查询用户头像
	headPath, err := h.userServicer.GetUserHeadPath(ctx, userName)
	if err != nil {
		logger.Error("查询用户头像失败", logger.S("user_name", userName), logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户头像失败"})
		return
	}
	// 返回用户头像
	c.JSON(http.StatusOK, gin.H{"head_path": headPath})
}

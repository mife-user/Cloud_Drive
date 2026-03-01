package handlers

import (
	"context"
	"drive/internal/api/dtos/request"
	"drive/internal/api/dtos/response"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	logger.Info("开始处理用户注册请求")
	defer logger.Info("用户注册请求处理完成")
	// 设置合理的超时时间，注册操作涉及数据库写入和缓存
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	var userAuthDto request.UserAuthDT
	if err := c.ShouldBindJSON(&userAuthDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userDOM := userAuthDto.ToDMUserAuth()
	err := h.userServicer.Register(ctx, userDOM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	// 转换为DTO
	userRegResp := response.ToDTUserReg(userDOM)
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data":    userRegResp,
	})
}

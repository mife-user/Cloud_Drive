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

// Login 用户登录
func (s *UserHandler) Login(c *gin.Context) {
	logger.Info("开始处理用户登录请求")
	defer logger.Info("用户登录请求处理完成")
	// 设置合理的超时时间，登录操作涉及数据库查询和缓存
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	var userAuthDto request.UserAuthDT
	if err := c.ShouldBindJSON(&userAuthDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}
	userDOM := userAuthDto.ToDMUserAuth()
	token, err := s.userServicer.Login(ctx, userDOM)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	// 转换为DTO
	userLoginRPS := response.ToDTUserLogin(userDOM, token)
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    userLoginRPS,
	})
}

package handlers

import (
	"drive/internal/api/dtos/request"
	"drive/internal/api/dtos/response"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
func (s *UserHandler) Login(c *gin.Context) {
	logger.Info("开始处理用户登录请求")
	defer logger.Info("用户登录请求处理完成")

	var userAuthDto request.UserAuthDT
	if err := c.ShouldBindJSON(&userAuthDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}
	userDOM := userAuthDto.ToDMUserAuth()
	token, err := s.userServicer.Login(c.Request.Context(), userDOM)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	userLoginRPS := response.ToDTUserLogin(userDOM, token)
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    userLoginRPS,
	})
}

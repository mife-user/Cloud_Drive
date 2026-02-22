package handlers

import (
	"context"
	"drive/internal/api/dtos"
	"drive/internal/domain"
	"drive/pkg/auth"
	"drive/pkg/conf"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userRepo domain.UserRepo
	config   *conf.Config
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(userRepo domain.UserRepo, config *conf.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		config:   config,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	logger.Info("开始处理用户登录请求")
	defer logger.Info("用户登录请求处理完成")
	// 设置合理的超时时间，登录操作涉及数据库查询和缓存
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	var userDtos dtos.UserDtos
	if err := c.ShouldBindJSON(&userDtos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}

	// 检查用户名是否为空
	if userDtos.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名为空"})
		return
	}

	user := userDtos.ToDMUser()
	if err := h.userRepo.Logon(ctx, user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	token, err := auth.GenerateToken(user.ID, user.Role, user.UserName, h.config.JWT.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"username": user.UserName,
			"role":     user.Role,
		},
		"token": token,
	})
}

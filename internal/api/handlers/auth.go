package handlers

import (
	"net/http"

	"drive/internal/domain"
	"drive/pkg/auth"
	"drive/pkg/conf"

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
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.userRepo.Logon(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	token, err := auth.GenerateToken(user.ID, user.UserName, h.config.JWT.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"role":     user.Role,
		},
		"token": token,
	})
}

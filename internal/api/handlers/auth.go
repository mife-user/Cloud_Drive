package handlers

import (
	"drive/internal/api/dtos"
	"drive/internal/domain"
	"drive/pkg/auth"
	"drive/pkg/conf"
	"net/http"

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
	var userDtos dtos.UserDtos
	if err := c.ShouldBindJSON(&userDtos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	user := userDtos.ToUser()
	if err := h.userRepo.Logon(c, user); err != nil {
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

			"username": user.UserName,
			"role":     user.Role,
		},
		"token": token,
	})
}

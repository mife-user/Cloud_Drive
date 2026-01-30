package handlers

import (
	"net/http"

	"drive/internal/domain"
	"drive/pkg/conf"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userRepo domain.UserRepo
	config   *conf.Config
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userRepo domain.UserRepo, config *conf.Config) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		config:   config,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.userRepo.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	// 这里可以实现获取用户信息的逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户信息成功",
		"user_id": userID,
	})
}

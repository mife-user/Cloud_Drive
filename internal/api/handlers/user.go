package handlers

import (
	"net/http"

	"drive/internal/api/dtos"
	"drive/internal/domain"
	"drive/internal/service"
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
	var userDto dtos.UserDtos
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if err := h.userRepo.Register(c, userDto.ToDMUser()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// 用户信息修改
func (h *UserHandler) RemixUser(c *gin.Context) {
	var userDto dtos.UserDtos
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	DMUser := userDto.ToDMUser()
	// 从上下文获取当前登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	if err := service.BuildUser(DMUser, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
		return
	}
	if err := h.userRepo.RemixUser(c, DMUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}

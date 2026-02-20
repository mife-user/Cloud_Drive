package handlers

import (
	"drive/internal/api/dtos"
	"drive/internal/service"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户信息修改
func (h *UserHandler) RemixUser(c *gin.Context) {
	logger.Info("开始处理用户信息修改请求")
	defer logger.Info("用户信息修改请求处理完成")
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
	// 修改用户信息
	if err := h.userRepo.RemixUser(c, DMUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}

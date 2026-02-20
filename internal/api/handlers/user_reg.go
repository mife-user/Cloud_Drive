package handlers

import (
	"drive/internal/api/dtos"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	logger.Info("开始处理用户注册请求")
	defer logger.Info("用户注册请求处理完成")
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

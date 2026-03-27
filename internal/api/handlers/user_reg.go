package handlers

import (
	"drive/internal/api/dtos/request"
	"drive/internal/api/dtos/response"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	logger.Info("开始处理用户注册请求")
	defer logger.Info("用户注册请求处理完成")

	var userAuthDto request.UserAuthDT
	if err := c.ShouldBindJSON(&userAuthDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userDOM := userAuthDto.ToDMUserAuth()
	err := h.userServicer.Register(c.Request.Context(), userDOM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	userRegResp := response.ToDTUserReg(userDOM)
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"data":    userRegResp,
	})
}

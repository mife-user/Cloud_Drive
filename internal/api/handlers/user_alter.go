package handlers

import (
	"drive/internal/api/dtos/request"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户信息修改
func (h *UserHandler) RemixUser(c *gin.Context) {
	logger.Info("开始处理用户信息修改请求")
	defer logger.Info("用户信息修改请求处理完成")

	var userHeaderDto request.UserAlterDT
	if err := c.ShouldBindJSON(&userHeaderDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	DMUser := userHeaderDto.ToDMUserAlter()

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	DMUser.ID = userIDUint

	if err := h.userServicer.RemixUser(c.Request.Context(), DMUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}

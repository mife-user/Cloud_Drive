package handlers

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) UpdateHeader(c *gin.Context) {
	logger.Info("开始处理更新用户头像请求")
	defer logger.Info("更新用户头像请求处理完成")
	ctx, cancel := context.WithTimeout(c, 1*time.Minute)
	defer cancel()
	// 获取用户名
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	usernameSTR, ok := exc.IsString(username)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户名"})
		return
	}
	// 获取用户上传的文件
	file, err := c.FormFile("header")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	// 获取用户角色
	userRole, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	// 检查用户角色是否为普通用户
	userRoleSTR, ok := exc.IsString(userRole)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "role类型错误"})
		return
	}
	fileRecord := &domain.UserHeader{
		Username:   usernameSTR,
		HeaderPath: "",
		Role:       userRoleSTR,
	}
	if err := h.userServicer.UpdateHeader(ctx, fileRecord, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户头像失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "用户头像更新成功"})
}

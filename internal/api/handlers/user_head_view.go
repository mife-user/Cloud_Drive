package handlers

import (
	"drive/internal/api/dtos/response"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetHeader(c *gin.Context) {
	logger.Info("开始处理用户信息修改请求")
	defer logger.Info("用户信息修改请求处理完成")

	username := c.Param("username")
	if username == "" {
		logger.Debug("获取用户信息失败" + errorer.ErrUserIDNotFound)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorer.ErrUserIDNotFound})
		return
	}

	headPath, err := h.userServicer.GetUserHeadPath(c.Request.Context(), username)
	if err != nil {
		logger.Error("查询用户头像失败", logger.S("username", username), logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户头像失败"})
		return
	}
	headPathRSP := response.ToDTHeaderPath(headPath)

	c.JSON(http.StatusOK, gin.H{
		"message": "查询用户头像成功",
		"data":    headPathRSP,
	})
}

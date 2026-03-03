package handlers

import (
	"context"
	"drive/internal/api/dtos/request"
	"drive/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) UploadFile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Minute)
	defer cancel()

	var fileDto request.FilePermissionsDT
	fileDto.Permissions = "public"

	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
		return
	}

	files := file.File["files"]

	userID, existsID := c.Get("user_id")
	if !existsID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userName, existsNM := c.Get("user_name")
	if !existsNM {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userRole, existsRole := c.Get("role")
	if !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, userNameSTR, userRoleSTR, ok := service.ExchangeType(userID, userName, userRole)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型转换失败"})
		return
	}

	totalSize := service.GetTotalSize(files)

	nowSize, ok := h.fileServicer.CheckUserSize(ctx, userIDUint, totalSize)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户空间不足"})
		return
	}

	fileKey := fileDto.ToDMFilePermissions()
	fileKey.UserID = userIDUint
	fileKey.Owner = userNameSTR

	fileRecords, err := service.SaveFiles(ctx, files, userRoleSTR, fileKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	if len(*fileRecords) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: 没有文件记录"})
		return
	}

	if err := h.fileServicer.UploadFile(ctx, *fileRecords, nowSize); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
	})
}

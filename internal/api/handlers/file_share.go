package handlers

import (
	"drive/internal/api/dtos"
	"drive/internal/service"
	"drive/pkg/errorer"
	"drive/pkg/exc"
	"drive/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShareFile 分享文件
func (h *FileHandler) ShareFile(c *gin.Context) {
	logger.Info("开始处理分享文件请求")

	var req dtos.ShareFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("绑定请求参数失败", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, existsID := c.Get("user_id")
	if !existsID {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userName, existsNM := c.Get("user_name")
	if !existsNM {
		logger.Error("获取用户名失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	// 转换userID和userName为uint和string类型
	userIDUint, userNameStr, err := service.ExchangeFile(userID, userName)
	if err != nil {
		logger.Error("转换用户ID和用户名失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 分享文件
	shareID, accessKey, err := h.fileRepo.ShareFile(c, req.FileID, userIDUint, userNameStr)
	if err != nil {
		logger.Error("分享文件失败", logger.S("file_id", fmt.Sprintf("%d", req.FileID)), logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 构建分享URL
	shareURL := fmt.Sprintf("/api/files/share/%s", shareID)

	logger.Info("分享文件成功", logger.S("share_id", shareID), logger.S("file_id", fmt.Sprintf("%d", req.FileID)))
	c.JSON(http.StatusOK, gin.H{
		"message":    "分享成功",
		"share_url":  shareURL,
		"access_key": accessKey,
	})
}

// AccessShare 访问文件分享
func (h *FileHandler) AccessShare(c *gin.Context) {
	logger.Info("开始处理访问分享请求")

	shareID := c.Param("share_id")
	if shareID == "" {
		logger.Error("分享ID为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "分享ID不能为空"})
		return
	}

	var req dtos.AccessShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("绑定请求参数失败", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	file, err := h.fileRepo.AccessShare(c.Request.Context(), shareID, req.AccessKey)
	if err != nil {
		logger.Error("访问分享失败", logger.S("share_id", shareID), logger.C(err))
		switch err.Error() {
		case errorer.ErrShareNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "分享不存在"})
		case errorer.ErrInvalidAccessKey:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "访问密钥无效"})
		case errorer.ErrShareExpired:
			c.JSON(http.StatusGone, gin.H{"error": "分享已过期"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "访问失败"})
		}
		return
	}

	logger.Info("访问分享成功", logger.S("share_id", shareID), logger.S("file_path", file.Path))
	c.File(file.Path)
}

// UpdateFilePermissions 更新文件权限
func (h *FileHandler) UpdateFilePermissions(c *gin.Context) {
	logger.Info("开始处理更新文件权限请求")

	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		logger.Error("文件ID为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	fileID, err := exc.StrToUint(fileIDStr)
	if err != nil {
		logger.Error("文件ID格式错误", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID格式错误"})
		return
	}

	var req dtos.FileDtos
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("绑定请求参数失败", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, existsID := c.Get("user_id")
	if !existsID {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	if err := h.fileRepo.UpdateFilePermissions(c.Request.Context(), fileID, userIDUint, req.Permissions); err != nil {
		logger.Error("更新文件权限失败", logger.S("file_id", fileIDStr), logger.C(err))
		switch err.Error() {
		case errorer.ErrFileNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		case errorer.ErrNotFileOwner:
			c.JSON(http.StatusForbidden, gin.H{"error": "非文件所有者，无法更新权限"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		}
		return
	}

	logger.Info("更新文件权限成功", logger.S("file_id", fileIDStr))
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
	})
}

// AddFavorite 添加文件收藏
func (h *FileHandler) AddFavorite(c *gin.Context) {
	logger.Info("开始处理添加收藏请求")

	var req dtos.FavoriteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("绑定请求参数失败", logger.C(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, existsID := c.Get("user_id")
	if !existsID {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	if err := h.fileRepo.AddFavorite(c.Request.Context(), userIDUint, req.FileID); err != nil {
		logger.Error("添加收藏失败", logger.S("file_id", fmt.Sprintf("%d", req.FileID)), logger.C(err))
		switch err.Error() {
		case errorer.ErrFileNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		case errorer.ErrNotFileOwner:
			c.JSON(http.StatusForbidden, gin.H{"error": "非文件所有者，无法收藏"})
		case errorer.ErrFavoriteExist:
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件已收藏"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加收藏失败"})
		}
		return
	}

	logger.Info("添加收藏成功", logger.S("file_id", fmt.Sprintf("%d", req.FileID)))
	c.JSON(http.StatusOK, gin.H{
		"message": "收藏成功",
	})
}

// RemoveFavorite 取消文件收藏
func (h *FileHandler) RemoveFavorite(c *gin.Context) {
	logger.Info("开始处理取消收藏请求")

	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	fileID, err := exc.StrToUint(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID格式错误"})
		return
	}

	userID, existsID := c.Get("user_id")
	if !existsID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	if err := h.fileRepo.RemoveFavorite(c.Request.Context(), userIDUint, fileID); err != nil {

		switch err.Error() {
		case errorer.ErrFavoriteNotExist:
			c.JSON(http.StatusNotFound, gin.H{"error": "收藏不存在"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		}
		return
	}

	logger.Info("取消收藏成功", logger.S("file_id", fileIDStr))
	c.JSON(http.StatusOK, gin.H{
		"message": "取消收藏成功",
	})
}

// GetFavorites 获取收藏列表
func (h *FileHandler) GetFavorites(c *gin.Context) {
	logger.Info("开始处理获取收藏列表请求")

	userID, existsID := c.Get("user_id")
	if !existsID {
		logger.Error("获取用户ID失败")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("用户ID类型转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "类型错误"})
		return
	}

	files, err := h.fileRepo.GetFavorites(c.Request.Context(), userIDUint)
	if err != nil {
		logger.Error("获取收藏列表失败", logger.C(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败"})
		return
	}

	logger.Info("获取收藏列表成功", logger.S("count", fmt.Sprintf("%d", len(files))))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"files":   files,
	})
}

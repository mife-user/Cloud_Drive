package handlers

import (
	"drive/internal/api/dtos/request"
	"drive/internal/api/dtos/response"
	"drive/pkg/exc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) ShareFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID类型错误"})
		return
	}

	userName, exists := c.Get("user_name")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userNameSTR, ok := exc.IsString(userName)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名类型错误"})
		return
	}

	var req request.ShareFileDT
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	shareID, accessKey, err := h.fileServicer.ShareFile(c.Request.Context(), req.FileID, userIDUint, userNameSTR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ToDTFileShare(shareID, accessKey))
}

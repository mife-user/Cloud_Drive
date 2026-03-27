package handlers

import (
	"drive/internal/api/dtos/response"
	"drive/pkg/exc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) ViewFile(c *gin.Context) {
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

	file, err := h.fileServicer.ViewFile(c.Request.Context(), fileID, userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ToDTFileInfo(file))
}

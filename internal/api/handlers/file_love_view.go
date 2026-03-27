package handlers

import (
	"drive/internal/api/dtos/response"
	"drive/pkg/exc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) GetFavorites(c *gin.Context) {
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

	files, err := h.fileServicer.GetFavorites(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ToDTFavoriteList(files))
}

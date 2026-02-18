package middlewares

import (
	"net/http"
	"path/filepath"

	"drive/pkg/conf"

	"github.com/gin-gonic/gin"
)

// TypeCheck 文件类型检查中间件
func TypeCheck(config *conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
			c.Abort()
			return
		}

		files := form.File["files"]
		for _, fileHeader := range files {
			ext := filepath.Ext(fileHeader.Filename)
			if !config.Upload.TellType(ext) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "文件类型不允许: " + ext})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

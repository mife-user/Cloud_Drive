package handlers

import (
	"drive/internal/api/dtos"
	"drive/internal/domain"
	"drive/internal/service"
	"drive/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) UploadFile(c *gin.Context) {
	logger.Info("开始处理文件上传请求")
	defer logger.Info("文件上传请求处理完成")

	var fileDto dtos.FileDtos
	if err := c.ShouldBind(&fileDto); err != nil {
	}

	file, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
		return
	}

	files := file.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有文件"})
		return
	}

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

	userRole, existsRole := c.Get("user_role")
	if !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}

	userIDUint, _ := userID.(uint)
	userNameStr, _ := userName.(string)

	if fileDto.IsChunked {
		result, err := service.HandleChunkUpload(c, h.fileRepo, files[0], &fileDto, userIDUint, userNameStr, fileDto.Permissions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "分片上传失败: " + err.Error()})
			return
		}

		if result.Completed {
			fileRecords := []*domain.File{result.FileRecord}
			if err := h.fileRepo.UploadFile(c, fileRecords); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: " + err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message":        "上传成功",
				"upload_task_id": result.UploadTaskID,
				"completed":      true,
				"progress":       100,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":        "分片上传成功",
				"upload_task_id": result.UploadTaskID,
				"completed":      false,
				"progress":       result.Progress,
			})
		}
		return
	}

	fileRecords, err := service.SaveFiles(files, userID, userName, userRole, fileDto.ToDMFile())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败: " + err.Error()})
		return
	}

	if len(*fileRecords) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: 没有文件记录"})
		return
	}
	if err := h.fileRepo.UploadFile(c, *fileRecords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
	})
}

func (h *FileHandler) GetUploadStatus(c *gin.Context) {
	logger.Info("开始查询上传状态")
	defer logger.Info("上传状态查询完成")

	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	status, err := service.GetUploadStatus(c, h.fileRepo, uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询上传状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"status":  status,
	})
}

func (h *FileHandler) ViewFiles(c *gin.Context) {
	logger.Info("开始处理查看所有文件请求")
	defer logger.Info("查看所有文件请求处理完成")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证用户"})
		return
	}
	files, err := h.fileRepo.ViewFile(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查看文件失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查看成功",
		"files":   files,
	})
}

package service

import (
	"context"
	"drive/internal/api/dtos"
	"drive/internal/domain"
	"drive/pkg/errorer"
	"drive/pkg/logger"
	"drive/pkg/pool"
	"drive/pkg/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type ChunkUploadResult struct {
	UploadTaskID uint
	Completed    bool
	FileRecord   *domain.File
	Progress     float64
}

func SaveFiles(files []*multipart.FileHeader, userID any, userName any, userRole any, filekey *domain.File) (*[]*domain.File, error) {
	userRoleStr, ok := userRole.(string)
	if !ok {
		logger.Error("userRole类型转换失败")
		return nil, errorer.New(errorer.ErrTypeError)
	}
	if userRoleStr != "VIP" {
		filekey.Size = 1024 * 1024 * 1024
	} else {
		filekey.Size = 1024 * 1024 * 2048
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("userID类型转换失败")
		return nil, errorer.New(errorer.ErrTypeError)
	}
	userNameStr, ok := userName.(string)
	if !ok {
		logger.Error("userName类型转换失败")
		return nil, errorer.New(errorer.ErrTypeError)
	}
	filekey.Owner = userNameStr
	filekey.UserID = userIDUint
	recordCh := make(chan *domain.File, len(files))
	pool := pool.NewPool(4)
	pool.Start()
	var wg sync.WaitGroup
	for _, header := range files {
		h := header
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			fileRecord, err := utils.SaveFile(h, filekey)
			if err != nil {
				logger.Error("保存文件失败: %v", logger.C(err))
				return
			}
			recordCh <- fileRecord
		})
	}
	wg.Wait()
	close(recordCh)
	fileRecords := make([]*domain.File, 0, len(files))
	for record := range recordCh {
		fileRecords = append(fileRecords, record)
	}
	pool.Stop()
	return &fileRecords, nil
}

func HandleChunkUpload(ctx context.Context, fileRepo domain.FileRepo, header *multipart.FileHeader, fileDto *dtos.FileDtos, userID uint, userName string, permissions string) (*ChunkUploadResult, error) {
	var task *domain.UploadTask
	var err error

	if fileDto.UploadTaskID == 0 {
		task, err = fileRepo.GetUploadTaskByMD5(ctx, userID, fileDto.FileMD5)
		if err == nil && task != nil {
			fileDto.UploadTaskID = task.ID
		} else {
			task = &domain.UploadTask{
				UserID:          userID,
				FileName:        header.Filename,
				FileSize:        header.Size,
				FileMD5:         fileDto.FileMD5,
				TotalChunks:     fileDto.TotalChunks,
				CompletedChunks: "[]",
				Status:          0,
			}
			if err := fileRepo.CreateUploadTask(ctx, task); err != nil {
				return nil, err
			}
			fileDto.UploadTaskID = task.ID
		}
	} else {
		task, err = fileRepo.GetUploadTaskByID(ctx, fileDto.UploadTaskID)
		if err != nil {
			return nil, err
		}
	}

	if err := utils.SaveChunk(header, userID, fileDto.UploadTaskID, fileDto.ChunkIndex); err != nil {
		return nil, err
	}

	completedChunks, _ := utils.ParseCompletedChunks(task.CompletedChunks)
	chunkExists := false
	for _, idx := range completedChunks {
		if idx == fileDto.ChunkIndex {
			chunkExists = true
			break
		}
	}
	if !chunkExists {
		completedChunks = append(completedChunks, fileDto.ChunkIndex)
		task.CompletedChunks = utils.SerializeCompletedChunks(completedChunks)
		if err := fileRepo.UpdateUploadTask(ctx, task); err != nil {
			return nil, err
		}
	}

	progress := float64(len(completedChunks)) / float64(task.TotalChunks) * 100

	if len(completedChunks) == task.TotalChunks {
		storageBase := "./storage/" + strconv.FormatUint(uint64(userID), 10)
		if err := os.MkdirAll(storageBase, 0755); err != nil {
			return nil, err
		}

		ext := filepath.Ext(task.FileName)
		baseName := task.FileName[:len(task.FileName)-len(ext)]
		finalFileName := task.FileName
		tempPath := filepath.Join(storageBase, finalFileName)

		if _, err := os.Stat(tempPath); err == nil {
			finalFileName = baseName + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext
			tempPath = filepath.Join(storageBase, finalFileName)
		}

		if err := utils.MergeChunks(userID, fileDto.UploadTaskID, task.TotalChunks, tempPath); err != nil {
			return nil, err
		}

		calculatedMD5, err := utils.CalculateFileMD5(tempPath)
		if err != nil {
			return nil, err
		}
		if calculatedMD5 != task.FileMD5 {
			os.Remove(tempPath)
			utils.CleanupChunks(userID, fileDto.UploadTaskID)
			return nil, errorer.New("文件校验失败")
		}

		task.Status = 1
		fileRepo.UpdateUploadTask(ctx, task)

		fileRecord := &domain.File{
			FileName:    task.FileName,
			Size:        task.FileSize,
			Path:        tempPath,
			UserID:      userID,
			Owner:       userName,
			Permissions: permissions,
		}

		utils.CleanupChunks(userID, fileDto.UploadTaskID)

		return &ChunkUploadResult{
			UploadTaskID: fileDto.UploadTaskID,
			Completed:    true,
			FileRecord:   fileRecord,
			Progress:     100,
		}, nil
	}

	return &ChunkUploadResult{
		UploadTaskID: fileDto.UploadTaskID,
		Completed:    false,
		Progress:     progress,
	}, nil
}

func GetUploadStatus(ctx context.Context, fileRepo domain.FileRepo, taskID uint) (*dtos.UploadStatusResponse, error) {
	task, err := fileRepo.GetUploadTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	completedChunks, _ := utils.ParseCompletedChunks(task.CompletedChunks)
	progress := float64(len(completedChunks)) / float64(task.TotalChunks) * 100

	return &dtos.UploadStatusResponse{
		UploadTaskID:     task.ID,
		FileName:         task.FileName,
		FileSize:         task.FileSize,
		TotalChunks:      task.TotalChunks,
		CompletedChunks:  completedChunks,
		Progress:         progress,
		Status:           task.Status,
	}, nil
}

func ExchangeFile(userID any, userName any) (uint, string, error) {
	userIDUint, ok := userID.(uint)
	if !ok {
		logger.Error("userID类型转换失败")
		return 0, "", errorer.New(errorer.ErrTypeError)
	}
	userNameStr, ok := userName.(string)
	if !ok {
		logger.Error("userName类型转换失败")
		return 0, "", errorer.New(errorer.ErrTypeError)
	}
	return userIDUint, userNameStr, nil
}

func ParseID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Error("ID类型转换失败")
		return 0, errorer.New(errorer.ErrTypeError)
	}
	return uint(id), nil
}

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const DefaultChunkSize = 5 * 1024 * 1024

func GetChunkDir(userID uint, uploadTaskID uint) string {
	return fmt.Sprintf("./storage/chunks/%d/%d", userID, uploadTaskID)
}

func SaveChunk(header *multipart.FileHeader, userID uint, uploadTaskID uint, chunkIndex int) error {
	chunkDir := GetChunkDir(userID, uploadTaskID)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return fmt.Errorf("创建分片目录失败: %w", err)
	}

	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", chunkIndex))

	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("打开分片文件失败: %w", err)
	}
	defer file.Close()

	dst, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("创建分片文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("保存分片文件失败: %w", err)
	}

	return nil
}

func MergeChunks(userID uint, uploadTaskID uint, totalChunks int, destPath string) error {
	chunkDir := GetChunkDir(userID, uploadTaskID)

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("打开分片 %d 失败: %w", i, err)
		}

		if _, err := io.Copy(dst, chunkFile); err != nil {
			chunkFile.Close()
			return fmt.Errorf("合并分片 %d 失败: %w", i, err)
		}
		chunkFile.Close()
	}

	return nil
}

func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ParseCompletedChunks(completedStr string) ([]int, error) {
	if completedStr == "" {
		return []int{}, nil
	}

	var chunks []int
	err := json.Unmarshal([]byte(completedStr), &chunks)
	if err != nil {
		parts := strings.Split(completedStr, ",")
		chunks = make([]int, 0, len(parts))
		for _, part := range parts {
			if idx, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				chunks = append(chunks, idx)
			}
		}
	}
	return chunks, nil
}

func SerializeCompletedChunks(chunks []int) string {
	data, _ := json.Marshal(chunks)
	return string(data)
}

func CleanupChunks(userID uint, uploadTaskID uint) error {
	chunkDir := GetChunkDir(userID, uploadTaskID)
	return os.RemoveAll(chunkDir)
}

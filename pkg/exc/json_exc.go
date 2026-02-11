package exc

import (
	"drive/internal/domain"
	"drive/pkg/logger"
	"encoding/json"
)

// 将domain.File序列化为JSON字符串
func ExcFileToJSON(file domain.File) (string, error) {
	fileJSON, err := json.Marshal(file)
	if err != nil {
		logger.Error("序列化文件失败", logger.C(err))
		return "", err
	}
	return string(fileJSON), nil
}

// 将JSON字符串反序列化为domain.File
func ExcJSONToFile(fileJSON string, file *domain.File) error {
	err := json.Unmarshal([]byte(fileJSON), file)
	if err != nil {
		logger.Error("反序列化文件失败", logger.C(err))
		return err
	}
	return nil
}

package service

import (
	"context"
	"drive/internal/domain"
	"drive/pkg/exc"
	"drive/pkg/save"
	"mime/multipart"
)

func ExchangeType(userID, userName, userRole any) (uint, string, string, bool) {
	userIDUint, ok := exc.IsUint(userID)
	if !ok {
		return 0, "", "", false
	}
	userNameSTR, ok := exc.IsString(userName)
	if !ok {
		return 0, "", "", false
	}
	userRoleSTR, ok := exc.IsString(userRole)
	if !ok {
		return 0, "", "", false
	}
	return userIDUint, userNameSTR, userRoleSTR, true
}

func GetTotalSize(files []*multipart.FileHeader) int64 {
	var total int64
	for _, file := range files {
		total += file.Size
	}
	return total
}

func SaveFiles(ctx context.Context, files []*multipart.FileHeader, role string, fileKey *domain.File) (*[]*domain.File, error) {
	var fileRecords []*domain.File
	var rating int
	if role == "vip" {
		rating = 1024 * 1024
	} else {
		rating = 512 * 1024
	}
	for _, fileHeader := range files {
		fileName, size, path, err := save.SaveFile(ctx, fileHeader, 10*1024*1024, fileKey.UserID, rating)
		if err != nil {
			return nil, err
		}
		fileRecords = append(fileRecords, &domain.File{
			FileName:    fileName,
			Size:        size,
			Path:        path,
			UserID:      fileKey.UserID,
			Owner:       fileKey.Owner,
			Permissions: fileKey.Permissions,
		})
	}
	return &fileRecords, nil
}

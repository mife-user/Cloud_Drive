package service

import (
	"context"
)

// ShareFile 分享文件
func (s *fileServicer) ShareFile(ctx context.Context, fileID uint, userID uint, owner string) (string, string, error) {
	var err error
	var shareID, accessKey string
	if shareID, accessKey, err = s.fileRepo.ShareFile(ctx, fileID, userID, owner); err != nil {
		return "", "", err
	}
	return shareID, accessKey, nil
}

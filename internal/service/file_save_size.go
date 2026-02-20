package service

import "mime/multipart"

//获取上传文件所有大小
func GetTotalSize(files []*multipart.FileHeader) int64 {
	var totalSize int64
	for _, file := range files {
		totalSize += file.Size
	}
	return totalSize
}

package conf

import "strings"

// initAllowedTypesSet 初始化允许文件类型的 Set 用于快速查找
func (u *UploadConfig) InitAllowedTypesSet() {
	u.allowedTypesSet = make(map[string]bool)
	for _, ext := range u.AllowedFileTypes {
		u.allowedTypesSet[strings.ToLower(ext)] = true
	}
}

// TellType 检查文件类型是否允许
func (u *UploadConfig) TellType(ext string) bool {
	return u.allowedTypesSet[strings.ToLower(ext)]
}

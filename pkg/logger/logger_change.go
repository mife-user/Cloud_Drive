package logger

import "go.uber.org/zap"

// 转换错误为zap.Field
func C(err error) zap.Field {
	return zap.Error(err)
}

// 转换字符串为zap.Field
func S(key string, value string) zap.Field {
	return zap.String(key, value)
}

package logger

import "go.uber.org/zap"

// logger 日志结构体
type logger struct {
	Log *zap.Logger
	env string
}

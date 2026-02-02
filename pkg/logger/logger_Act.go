package logger

import "go.uber.org/zap"

// Info 打印info日志
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

// Error 打印error日志
func (l *logger) Error(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}

// Debug 打印debug日志
func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.Log.Debug(msg, fields...)
}

// Warn 打印warn日志
func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

// Fatal 打印fatal日志
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.Log.Fatal(msg, fields...)
}

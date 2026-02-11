package logger

import (
	"drive/pkg/conf"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerInstance *zap.Logger

// 初始化日志实例
func InitLogger(config *conf.Config) error {
	// 日志目录
	logDir := "./logs"
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}
	// 日志文件路径
	logFile := filepath.Join(logDir, "app.log")
	// 打开日志文件
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// 获取当前环境配置
	env := config.Env
	// 日志编码器配置
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "time",       // 日志时间键名
		LevelKey:      "level",      // 日志级别键名
		NameKey:       "logger",     // 日志实例名键名
		CallerKey:     "caller",     // 日志调用者键名
		MessageKey:    "msg",        // 日志消息键名
		StacktraceKey: "stacktrace", // 日志栈跟踪键名
		LineEnding:    "||\n",       // 自定义日志行结束符
		EncodeTime:    nowTime,
	}

	// 根据环境选择编码器
	var encoder zapcore.Encoder
	if env == "prod" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	// 文件日志写入器
	fileWrite := zapcore.Lock(file)
	// 控制台日志写入器
	consoleWrite := zapcore.Lock(os.Stdout)

	// 根据环境选择输出目标
	var multiWrite zapcore.WriteSyncer
	if env == "prod" {
		multiWrite = zapcore.NewMultiWriteSyncer(fileWrite)
	} else {
		multiWrite = zapcore.NewMultiWriteSyncer(fileWrite, consoleWrite)
	}

	// 根据环境选择日志级别
	level := zapcore.DebugLevel // 日志级别
	if env == "prod" {
		level = zapcore.InfoLevel
	}
	// 日志核心
	core := zapcore.NewCore(encoder, multiWrite, level)
	// 日志实例
	loggerInstance = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

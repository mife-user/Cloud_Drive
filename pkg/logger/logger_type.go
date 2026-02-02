package logger

import "go.uber.org/zap"

type logger struct {
	Log *zap.Logger
	env string
}

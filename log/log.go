package log

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Level = logrus.DebugLevel
}

// Info 方法，记录一条 Info 级别的日志
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Error 方法，记录一条 Error 级别的日志
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Warn 方法，记录一条 Warn 级别的日志
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Debug 方法，记录一条 Debug 级别的日志
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

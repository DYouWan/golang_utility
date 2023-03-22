package log

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Level = logrus.DebugLevel
	// 设置输出格式为 TextFormatter
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
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

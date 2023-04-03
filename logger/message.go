package logger

import (
	"sync"
	"time"
)

var logMessagePool = sync.Pool{
	New: func() interface{} {
		return &LogMessage{}
	},
}

// LogMessage 日志消息结构体
type LogMessage struct {
	level  Level     // 日志等级
	time   time.Time // 日志时间
	msg    string    // 日志内容
	source string    // 日志来源（可选）
}

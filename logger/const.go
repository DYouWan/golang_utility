package logger

import "fmt"

const (
	defaultDir        = "./logs"
	defaultMaxSize    = int64(1024)
	defaultBufferSize = 1024
)

// Level 日志等级类型
type Level uint8

const (
	// PanicLevel 级别，最高级别的严重程度。记录日志并使用 Debug、Info、... 传递的消息调用 panic。
	PanicLevel Level = iota
	// FatalLevel 级别。记录日志并调用 `logger.Exit(1)`。即使日志级别设置为 Panic，它也将退出。
	FatalLevel
	// ErrorLevel 级别。记录日志。用于绝对需要注意的错误。通常用于钩子将错误发送到错误跟踪服务。
	ErrorLevel
	// WarnLevel 级别。非关键性条目，值得注意。
	WarnLevel
	// InfoLevel 级别。有关应用程序内部正在发生的操作的常规操作条目。
	InfoLevel
	// DebugLevel 级别。通常仅在调试时启用。非常详细的日志记录。
	DebugLevel
	// TraceLevel 级别。指定比 Debug 更细粒度的信息事件。
	TraceLevel
)

// AllLevels 公开了所有的日志记录级别
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

// String 将Level转换为字符串。
func (level Level) String() string {
	if b, err := level.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

func (level Level) MarshalText() ([]byte, error) {
	switch level {
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid logrus level %d", level)
}

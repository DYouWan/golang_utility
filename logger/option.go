package logger

var DefaultOptions = Options{
	Dir:        defaultDir,
	Level:      DebugLevel,
	MaxSize:    defaultMaxSize,
	BufferSize: defaultBufferSize,
}

// Options 日志选项
type Options struct {
	Dir        string // 日志文件目录
	MaxSize    int64  // 日志文件大小，单位字节
	BufferSize int    // 日志缓冲区大小，单位条
	Level      Level  // 日志级别
}

package log

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var (
	logDir            = "./logs"
	DefaultLog        = new(Logger)
	defaultMaxSize    = int64(1024 * 1024 * 1024)
	defaultBufferSize = 1024
)

// Logger 日志对象
type Logger struct {
	mu      sync.Mutex
	file    *os.File
	options Options
	logger  *logrus.Logger
	maxSize int64
	buffers map[logrus.Level]chan string
}

// Options 是日志选项
type Options struct {
	Path       string       // 日志文件名
	BufferSize int          // 日志缓冲区大小
	Level      logrus.Level // 日志级别
}

func init() {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	filePath := fmt.Sprintf("%s/%s.log", logDir, time.Now().Format("2006-01-02"))
	defaultOptions := Options{
		Path:       filePath,
		BufferSize: 1024,
		Level:      logrus.DebugLevel,
	}

	var err error
	DefaultLog, err = NewLogger(defaultOptions)
	if err != nil {
		panic(err)
	}
}

// NewLogger 创建一个新的日志对象
func NewLogger(options Options) (*Logger, error) {
	if options.BufferSize <= 0 {
		options.BufferSize = defaultBufferSize
	}

	file, err := os.OpenFile(options.Path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	buffers := make(map[logrus.Level]chan string)
	for _, level := range logrus.AllLevels {
		buffers[level] = make(chan string, options.BufferSize)
	}

	logger := &Logger{
		file:    file,
		options: options,
		buffers: buffers,
		maxSize: defaultMaxSize,
	}

	log := logrus.New()
	log.Level = options.Level

	// 设置输出格式为 TextFormatter
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	// 添加文件轮转钩子
	log.Hooks.Add(NewRotateHook(logger))

	logger.logger = log

	for level, buffer := range logger.buffers {
		go logger.processMessages(level, buffer)
	}

	return logger, nil
}

// Log 记录日志
func (log *Logger) Log(level logrus.Level, args ...interface{}) error {
	buffer, ok := log.buffers[level]
	if !ok {
		return errors.New("invalid log level")
	}
	log.mu.Lock()
	defer log.mu.Unlock()

	// 检查文件是否已关闭
	if log.file == nil {
		return fmt.Errorf("log file is closed")
	}

	entry := log.logger.WithField("source", "myapp")

	switch len(args) {
	case 1:
		entry = entry.WithField("message", args[0])
	case 2:
		entry = entry.WithFields(logrus.Fields{
			"message": args[0],
			"data":    args[1],
		})
	default:
		return fmt.Errorf("invalid number of arguments: %d", len(args))
	}

	select {
	case buffer <- entry.Message:
		return nil
	default:
		return errors.New("log buffer is full")
	}
}

// Close 关闭日志文件
func (log *Logger) Close() {
	if log.file != nil {
		for _, buffer := range log.buffers {
			for len(buffer) > 0 {
				msg := <-buffer
				fmt.Fprintln(log.file, msg)
			}
		}
		log.file.Close()
		log.file = nil
	}
}

// Info 方法，记录一条 Info 级别的日志
func Info(args ...interface{}) {
	DefaultLog.Log(logrus.InfoLevel, args)
}

// Error 方法，记录一条 Error 级别的日志
func Error(args ...interface{}) {
	DefaultLog.Log(logrus.ErrorLevel, args)
}

// Warn 方法，记录一条 Warn 级别的日志
func Warn(args ...interface{}) {
	DefaultLog.Log(logrus.WarnLevel, args)
}

// Debug 方法，记录一条 Debug 级别的日志
func Debug(args ...interface{}) {
	DefaultLog.Log(logrus.DebugLevel, args)
}

// 获取日志消息
func getMessage(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	switch msg := args[0].(type) {
	case string:
		return msg
	default:
		return ""
	}
}

func (log *Logger) processMessages(level logrus.Level, buffer chan string) {
	switch level {
	case logrus.DebugLevel:
		go func() {
			for {
				select {
				case msg := <-buffer:
					log.logger.Debug(msg)
				}
			}
		}()
	case logrus.InfoLevel:
		go func() {
			for {
				select {
				case msg := <-buffer:
					log.logger.Info(msg)
				}
			}
		}()
	case logrus.WarnLevel:
		go func() {
			for {
				select {
				case msg := <-buffer:
					log.logger.Warn(msg)
				}
			}
		}()
	case logrus.ErrorLevel:
		go func() {
			for {
				select {
				case msg := <-buffer:
					log.logger.Error(msg)
				}
			}
		}()
	case logrus.FatalLevel:
		go func() {
			for {
				select {
				case msg := <-buffer:
					log.logger.Fatal(msg)
				}
			}
		}()
	case logrus.PanicLevel:
		go func() {
			for {
				select {
				case msg := <-buffer:
					log.logger.Panic(msg)
				}
			}
		}()
	}
}

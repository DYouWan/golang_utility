package logger

import (
	"fmt"
	"github.com/dyouwan/utility/file"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

var (
	logDir            = "./logs"
	DefaultLog        = new(Logger)
	defaultMaxSize    = int64(1024)
	defaultBufferSize = 10
)

// Logger 日志对象
type Logger struct {
	mu      sync.Mutex
	options Options
	files   map[logrus.Level]*os.File
	buffers map[logrus.Level]chan string
	loggers map[logrus.Level]*logrus.Logger
}

// Options 是日志选项
type Options struct {
	MaxSize    int64        // 日志文件大小，单位字节
	BufferSize int          // 日志缓冲区大小，单位条
	Level      logrus.Level // 日志级别
}

func init() {
	defaultOptions := Options{
		MaxSize:    defaultMaxSize,
		BufferSize: defaultBufferSize,
		Level:      logrus.DebugLevel,
	}

	_, err := NewLogger(defaultOptions)
	if err != nil {
		panic(err)
	}
}

// NewLogger 创建一个新的日志对象
func NewLogger(options Options) (*Logger, error) {
	if options.MaxSize <= 0 {
		options.MaxSize = defaultMaxSize
	}

	if options.BufferSize <= 0 {
		options.BufferSize = defaultBufferSize
	}

	err := file.CrateFile(logDir)
	if err != nil {
		panic(err)
	}

	files := make(map[logrus.Level]*os.File)
	buffers := make(map[logrus.Level]chan string)
	loggers := make(map[logrus.Level]*logrus.Logger)
	for _, level := range logrus.AllLevels {
		// 判断文件夹是否存在，如果不存在就创建
		levelPath := fmt.Sprintf("%s/%s", logDir, level.String())
		err = file.CrateFile(levelPath)
		if err != nil {
			panic(err)
		}

		//根据AllLevels 初始化每个日志等级的缓冲区大小
		buffers[level] = make(chan string, options.BufferSize)

		//根据当前日期获取对应的文件,如果文件不存在就创建
		filePath := fmt.Sprintf("%s/%s.log", levelPath, time.Now().Format("2006-01-02"))
		var file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		files[level] = file
		loggers[level] = createLogger(level, levelPath, options.MaxSize, file)
	}

	logger := &Logger{
		options: options,
		files:   files,
		buffers: buffers,
		loggers: loggers,
	}

	for level, buffer := range logger.buffers {
		go logger.processMessages(level, buffer)
	}

	DefaultLog = logger
	return logger, nil
}

func (log *Logger) processMessages(level logrus.Level, buffer <-chan string) {
	for msg := range buffer {
		switch level {
		case logrus.DebugLevel:
			log.loggers[level].Debug(msg)
		case logrus.InfoLevel:
			log.loggers[level].Info(msg)
		case logrus.WarnLevel:
			log.loggers[level].Warn(msg)
		case logrus.ErrorLevel:
			log.loggers[level].Error(msg)
		case logrus.FatalLevel:
			log.loggers[level].Fatal(msg)
		case logrus.PanicLevel:
			log.loggers[level].Panic(msg)
		}
	}
}

func createLogger(level logrus.Level, levelPath string, maxSize int64, file *os.File) *logrus.Logger {
	log := logrus.New()
	log.Out = file
	log.Level = level
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	//// 添加文件轮转钩子
	//hook := &Hook{
	//	file:      file,
	//	maxSize:   maxSize,
	//	logger:    log,
	//	levelPath: levelPath,
	//}
	//log.AddHook(hook)

	return log
}

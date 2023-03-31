package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

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

// Log 记录日志
func (log *Logger) Log(level logrus.Level, args ...interface{}) error {
	file, ok := log.files[level]
	if !ok {
		return errors.New("invalid log level")
	}

	if file == nil {
		return fmt.Errorf("log file is closed")
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	err = log.IsFileNeedBackupByDate(file.Name(), level)
	if err != nil {
		return err
	}

	if fileInfo.Size() >= log.options.MaxSize {
		err = log.IsFileNeedBackupBySize(file.Name(), level)
		if err != nil {
			return err
		}
	}

	err = log.SendMessageToBuffer(level, args)
	if err != nil {
		return err
	}

	return nil
}

// IsFileNeedBackupByDate 根据日期判断是否需要备份文件
func (log *Logger) IsFileNeedBackupByDate(filePath string, level logrus.Level) error {
	today := time.Now().Format("2006-01-02")
	if filePath == "" {
		return errors.New("invalid file name")
	}

	fileName := path.Base(filePath)
	fileSuffix := filepath.Ext(fileName)
	filePrefix := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	if filePrefix == today {
		return nil
	}

	dirPath := path.Dir(filePath)
	newFilePath := fmt.Sprintf("%s/%s%s", dirPath, today, fileSuffix)

	log.mu.Lock()
	defer log.mu.Unlock()

	file, err := os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	log.files[level] = file
	log.loggers[level].Out = file

	return nil
}

// IsFileNeedBackupBySize 替换指定路径下的文件内容，并备份原文件。
func (log *Logger) IsFileNeedBackupBySize(filePath string, level logrus.Level) error {
	buffer, ok := log.buffers[level]
	if !ok {
		return errors.New("invalid log level")
	}
	for len(buffer) == 0 {
		fmt.Println("1111", len(buffer))
		log.Close(level)
		break
	}

	backupFilename := fmt.Sprintf("%s.%s", filePath, uuid.New().String())
	err := os.Rename(filePath, backupFilename)
	if err != nil {
		fmt.Println(222222, err)
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(121111)
		return err
	}

	log.files[level] = file
	log.loggers[level].Out = file

	return nil
}

// SendMessageToBuffer 向缓冲区发送消息。
func (log *Logger) SendMessageToBuffer(level logrus.Level, args ...interface{}) error {
	buffer, ok := log.buffers[level]
	if !ok {
		return errors.New("invalid log level")
	}

	entry := log.loggers[level].WithField("source", "myapp")
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

	msg, err := json.Marshal(entry.Data)
	if err != nil {
		return err
	}

	//当缓冲区满了，这里采用等待的方式。
	//如果使用default处理缓冲区满了的情况 需要考虑消息丢失的问题
	select {
	case buffer <- string(msg):
	}
	return nil
}

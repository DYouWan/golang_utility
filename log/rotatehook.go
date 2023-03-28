package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// RotateHook 日志轮转钩子
type RotateHook struct {
	log *Logger
}

func NewRotateHook(log *Logger) *RotateHook {
	return &RotateHook{log}
}

// Fire 在日志文件超过最大大小时触发轮转操作
func (hook *RotateHook) Fire(entry *logrus.Entry) error {
	fileInfo, err := os.Stat(hook.log.file.Name())
	if err != nil {
		return err
	}

	if fileInfo.Size() >= hook.log.maxSize {
		hook.log.mu.Lock()
		defer hook.log.mu.Unlock()

		// 先关闭文件
		hook.log.file.Close()

		// 备份旧文件
		backupFilename := fmt.Sprintf("%s.%d", hook.log.file.Name(), time.Now().Unix())
		err = os.Rename(hook.log.file.Name(), backupFilename)
		if err != nil {
			return err
		}

		// 创建新文件
		newFile, err := os.OpenFile(hook.log.file.Name(), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		// 更新 Logger 对象中的文件句柄和缓冲区
		hook.log.mu.Lock()
		defer hook.log.mu.Unlock()

		if hook.log.file != nil {
			hook.log.file.Close()
		}
		hook.log.file = newFile

		for level, buffer := range hook.log.buffers {
			close(buffer)
			hook.log.buffers[level] = make(chan string, hook.log.options.BufferSize)
			go hook.log.processMessages(level, hook.log.buffers[level])
		}
	}

	return nil
}

func (hook *RotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

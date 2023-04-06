package logger

import (
	"fmt"
)

// LogMessageJob 定义一个日志消息的处理器
type LogMessageJob struct {
	Index   int64
	message *LogMessage
}

func (j *LogMessageJob) Do() {
	file, ok := DefaultLog.files[j.message.level]
	if !ok {
		fmt.Println("Invalid log file level")
	}
	if file == nil {
		fmt.Println("File has been closed")
	}

	logMsg := fmt.Sprintf("[%s] %s %s\n", j.message.level.String(), j.message.time.Format("2006-01-02 15:04:05"), j.message.msg)
	_, err := file.Write([]byte(logMsg))
	if err != nil {
		fmt.Println("Failed to write log message:", err)
	} else {
		err = file.Sync() // 刷入磁盘
		if err != nil {
			fmt.Println("Failed to sync log writer:", err)
		}
	}
	// 将 LogMessage 对象归还给对象池
	logMessagePool.Put(j.message)
}

//// IsFileNeedBackupByDate 根据日期判断是否需要备份文件
//func (log *Logger) IsFileNeedBackupByDate(filePath string, level logrus.Level) error {
//	today := time.Now().Format("2006-01-02")
//	if filePath == "" {
//		return errors.New("invalid file name")
//	}
//
//	fileName := path.Base(filePath)
//	fileSuffix := filepath.Ext(fileName)
//	filePrefix := strings.TrimSuffix(fileName, filepath.Ext(fileName))
//	if filePrefix == today {
//		return nil
//	}
//
//	dirPath := path.Dir(filePath)
//	newFilePath := fmt.Sprintf("%s/%s%s", dirPath, today, fileSuffix)
//
//	log.mu.Lock()
//	defer log.mu.Unlock()
//
//	file, err := os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE, 0666)
//	if err != nil {
//		return err
//	}
//
//	log.files[level] = file
//	log.loggers[level].Out = file
//
//	return nil
//}
//
//// IsFileNeedBackupBySize 替换指定路径下的文件内容，并备份原文件。
//func (log *Logger) IsFileNeedBackupBySize(filePath string, level logrus.Level) error {
//	buffer, ok := log.buffers[level]
//	if !ok {
//		return errors.New("invalid log level")
//	}
//	for len(buffer) == 0 {
//		fmt.Println("1111", len(buffer))
//		log.Close(level)
//		break
//	}
//
//	backupFilename := fmt.Sprintf("%s.%s", filePath, uuid.New().String())
//	err := os.Rename(filePath, backupFilename)
//	if err != nil {
//		fmt.Println(222222, err)
//		return err
//	}
//
//	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
//	if err != nil {
//		fmt.Println(121111)
//		return err
//	}
//
//	log.files[level] = file
//	log.loggers[level].Out = file
//
//	return nil
//}

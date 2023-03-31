package logger

import (
	"github.com/sirupsen/logrus"
)

// CloseAll 关闭文件
func CloseAll() {
	DefaultLog.CloseAll()
}

// Close 关闭文件
func Close(level logrus.Level) {
	DefaultLog.Close(level)
}

// Close 关闭日志文件
func (log *Logger) Close(level logrus.Level) {
	//log.mu.Lock()
	//defer log.mu.Unlock()
	//
	file, ok := log.files[level]
	if !ok || file == nil {
		return
	}
	//buffer, ok := log.buffers[level]
	//if !ok {
	//	return
	//}
	//
	////如果缓冲区数据不为空，则将缓冲区数据一次性写入到当前文件中
	//bufferLen := len(buffer)
	//if bufferLen > 0 {
	//	newCh := make(chan string, bufferLen)
	//	for { // 复制 ch 中的所有元素到 newCh 中
	//		select {
	//		case s, ok := <-buffer:
	//			if !ok {
	//				// buffer 已经被关闭，退出循环
	//				return
	//			}
	//			select {
	//			case newCh <- s:
	//				// 成功将数据写入 newCh，继续循环
	//			default:
	//				// newCh 已经满了，退出循环
	//				return
	//			}
	//		default:
	//			// buffer 中没有数据可读，退出循环
	//			return
	//		}
	//	}
	//	close(newCh)
	//	//当一个 chan 被关闭后，使用 for range 循环读取 chan 时，当 chan 中的所有数据已经被读取完毕时，循环会自动退出
	//	log.processMessages(level, newCh)
	//}
	file.Close()
	log.files[level] = nil
}

// CloseAll 关闭日志文件
func (log *Logger) CloseAll() {
	for level, buffer := range log.buffers {
		file, ok := log.files[level]
		if !ok || file == nil {
			continue
		}
		if len(buffer) > 0 {
			log.processMessages(level, buffer)
		}
		file.Close()
		log.files[level] = nil
	}
}

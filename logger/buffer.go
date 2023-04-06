package logger

import (
	"sync"
)

// CircularBuffer 环形缓冲区结构体
type CircularBuffer struct {
	buffer         []*LogMessage // 缓冲区数组
	readIndex      int           // 当前读取位置
	writeIndex     int           // 当前写入位置
	size           int           // 缓冲区大小
	used           int           // 已用空间
	mutex          sync.Mutex    // 互斥锁，用于避免多个goroutine同时操作缓冲区
	writeSemaphore chan struct{} // 写信号量，用于控制写入操作
	readSemaphore  chan struct{} // 读信号量，用于控制读取操作
}

// NewCircularBuffer 创建一个新的环形缓冲区实例
func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buffer:         make([]*LogMessage, size),
		readIndex:      0,
		writeIndex:     0,
		size:           size,
		used:           0,
		writeSemaphore: make(chan struct{}, 1),
		readSemaphore:  make(chan struct{}, 1),
	}
}

// Write 向缓冲区中写入一条日志消息
func (c *CircularBuffer) Write(msg *LogMessage) {
	c.writeSemaphore <- struct{}{} // 获取写信号量，阻塞直到有足够空间写入日志消息
	c.mutex.Lock()                 // 获取互斥锁，避免其他goroutine同时访问缓冲区
	defer c.mutex.Unlock()
	if c.used == c.size { // 缓冲区已满，动态扩容
		newSize := c.size * 2
		newBuffer := make([]*LogMessage, newSize)
		for i := 0; i < c.used; i++ {
			newBuffer[i] = c.buffer[(c.readIndex+i)%c.size]
		}
		c.buffer = newBuffer
		c.readIndex = 0
		c.writeIndex = c.used
		c.size = newSize
	}

	c.buffer[c.writeIndex] = msg // 写入日志消息
	c.writeIndex = (c.writeIndex + 1) % c.size
	c.used++
	<-c.writeSemaphore // 释放写信号量
}

// Read 从缓冲区中读取一条日志消息
func (c *CircularBuffer) Read() *LogMessage {
	c.readSemaphore <- struct{}{} // 获取读信号量，阻塞直到有足够数据可读
	c.mutex.Lock()                // 获取互斥锁，避免其他goroutine同时访问缓冲区
	defer c.mutex.Unlock()
	if c.used == 0 { // 缓冲区为空，返回nil
		<-c.readSemaphore // 释放读信号量
		return nil
	}
	msg := c.buffer[c.readIndex] // 读取日志消息
	c.readIndex = (c.readIndex + 1) % c.size
	c.used--
	<-c.readSemaphore // 释放读信号量
	return msg
}

// WriteCircular 向缓冲区中写入一条日志消息,如果缓冲区满了 就覆盖
func (c *CircularBuffer) WriteCircular(msg *LogMessage) {
	c.writeSemaphore <- struct{}{} // 获取写信号量，阻塞直到有足够空间写入日志消息
	c.mutex.Lock()                 // 获取互斥锁，避免其他goroutine同时访问缓冲区
	defer c.mutex.Unlock()
	if c.used == c.size { // 缓冲区已满，覆盖最早的数据
		c.readIndex = (c.readIndex + 1) % c.size
	} else { // 缓冲区未满，更新已用空间
		c.used++
	}
	c.buffer[c.writeIndex] = msg // 写入日志消息
	c.writeIndex = (c.writeIndex + 1) % c.size
	<-c.writeSemaphore // 释放写信号量
}

package pipeline

import "net/http"

// Handler 管道处理程序
type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// HandlerFunc 定义实现Handler接口的函数
type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// ServeHTTP 实现Handler接口
func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

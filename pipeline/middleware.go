package pipeline

import "net/http"

// middleware 中间件
type middleware struct {
	handler Handler
	next    *middleware
}

// ServeHTTP 实现底层http.Handler接口
func (m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.next.ServeHTTP)
}

// VoidMiddleware 空的中间件，作为末尾使用
func VoidMiddleware() *middleware {
	return &middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		&middleware{},
	}
}

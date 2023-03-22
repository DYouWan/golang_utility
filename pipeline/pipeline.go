package pipeline

import (
	"fmt"
	"github.com/dyouwan/utility/log"
	"net/http"
	"time"
)

// Pipeline 负责处理http请求管道模型
type Pipeline struct {
	middleware middleware
	handlers   []Handler
}

// New 创建一个管道模型
func New(handlers ...Handler) *Pipeline {
	h := []Handler{timingMiddleware()}
	h = append(h, handlers...)

	return &Pipeline{
		handlers:   h,
		middleware: build(h),
	}
}

// ServeHTTP 实现底层http.Handler接口
func (p *Pipeline) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	p.middleware.ServeHTTP(rw, r)
}

func (p *Pipeline) Use(handler Handler) {
	if handler == nil {
		panic("handler cannot be nil")
	}

	p.handlers = append(p.handlers, handler)
	p.middleware = build(p.handlers)
}

// build 中间件链
func build(handlers []Handler) middleware {
	voidMiddleware := middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			rw.Write([]byte("voidMiddleware"))
		}),
		&middleware{},
	}

	switch len(handlers) {
	case 1:
		return middleware{handlers[0], &voidMiddleware}
	case 2:
		return middleware{handlers[0], &middleware{handlers[1], &voidMiddleware}}
	case 3:
		return middleware{handlers[0], &middleware{handlers[1], &middleware{handlers[2], &voidMiddleware}}}
	default:
		var stack []middleware
		for i := len(handlers) - 1; i >= 0; i-- {
			m := middleware{handlers[i], &voidMiddleware}
			if len(stack) > 0 {
				m.next = &stack[len(stack)-1]
			}
			stack = append(stack, m)
		}
		return stack[len(stack)-1]
	}
}

func timingMiddleware() HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// 记录请求开始时间
		start := time.Now()

		// 调用下一个处理器
		next(rw, r)

		// 计算请求处理时间
		duration := time.Since(start).Milliseconds()

		// 输出请求日志
		log.Info(fmt.Sprintf("%s request %s%s 耗时:%dms", r.RemoteAddr, r.Host, r.URL.String(), duration))
	}
}

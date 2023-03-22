package middleware

import (
	"fmt"
	"github.com/dyouwan/utility/log"
	"github.com/dyouwan/utility/pipeline"
	"net/http"
	"time"
)

func Time() pipeline.HandlerFunc {
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

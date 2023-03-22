package middleware

import (
	"fmt"
	"github.com/dyouwan/utility/log"
	"github.com/dyouwan/utility/pipeline"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Time() pipeline.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// 获取请求 ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将请求 ID 添加到响应头中
		rw.Header().Set("X-Request-ID", requestID)

		// 记录请求开始时间
		start := time.Now()

		// 调用下一个处理器
		next(rw, r)

		// 计算请求处理时间
		duration := time.Since(start).Milliseconds()

		// 输出请求日志
		log.Info(fmt.Sprintf("%s %s request %s%s 耗时:%dms", requestID, r.RemoteAddr, r.Host, r.URL.String(), duration))
	}
}

package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipMiddleware struct {
	level int
}

// NewGzip 返回gzip中间件
func NewGzip(level int) *GzipMiddleware {
	return &GzipMiddleware{level: level}
}

func (m *GzipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 检查客户端是否支持 gzip 压缩
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		next(w, r)
		return
	}

	// 设置响应头，表明响应内容已经被 gzip 压缩
	w.Header().Set("Content-Encoding", "gzip")

	// 创建 gzip 编码器，并将响应写入其中
	gz, err := gzip.NewWriterLevel(w, m.level)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer gz.Close()

	gzr := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
	next(gzr, r)
}

// 封装响应写入器
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查客户端是否支持 gzip 压缩
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// 设置响应头，表明响应内容已经被 gzip 压缩
		w.Header().Set("Content-Encoding", "gzip")

		// 创建 gzip 编码器，并将响应写入其中
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzr := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzr, r)
	})
}

func GzipMiddlewareWithLevel(level int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查客户端是否支持 gzip 压缩
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// 设置响应头，表明响应内容已经被 gzip 压缩
			w.Header().Set("Content-Encoding", "gzip")

			// 创建 gzip 编码器，并将响应写入其中
			gz, err := gzip.NewWriterLevel(w, level)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			gzr := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
			next.ServeHTTP(gzr, r)
		})
	}
}

// 封装响应写入器
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

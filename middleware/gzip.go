package middleware

import (
	"compress/gzip"
	"github.com/dyouwan/utility/pipeline"
	"io"
	"net/http"
	"strings"
)

func Gzip(level int) pipeline.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// 检查客户端是否支持 gzip 压缩
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(rw, r)
			return
		}

		// 设置响应头，表明响应内容已经被 gzip 压缩
		rw.Header().Set("Content-Encoding", "gzip")

		// 创建 gzip 编码器，并将响应写入其中
		gz, err := gzip.NewWriterLevel(rw, level)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		gzr := &gzipResponseWriter{Writer: gz, ResponseWriter: rw}
		next(gzr, r)
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

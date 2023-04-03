package logger

import (
	"strings"
	"testing"
)

func BenchmarkLog(b *testing.B) {
	str := strings.Repeat("a", 1024*1024)
	for i := 0; i < b.N; i++ {
		Debug(str, "test")
	}
}

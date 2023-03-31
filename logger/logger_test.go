package logger

import "testing"

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("hello, world!")
	}
}

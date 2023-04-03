package main

import (
	"fmt"
	"github.com/dyouwan/utility/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, _ := logger.NewLogger(logger.DefaultOptions)
	for i := 0; i < 100; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				logger.Debug(fmt.Sprintf("goroutine-%d message-%d", i, j), "test")
			}
		}(i)
	}
	// Ctrl+C 退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}

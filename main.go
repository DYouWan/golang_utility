package main

import (
	"fmt"
	"github.com/dyouwan/utility/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	for i := 0; i < 100; i++ {
		logger.Info(i, "qqq")
	}

	// Ctrl+C 退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}

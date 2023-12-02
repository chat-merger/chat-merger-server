package main

import (
	"chatmerger/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	app.Run()
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

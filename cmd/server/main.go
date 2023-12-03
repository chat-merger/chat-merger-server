package main

import (
	"chatmerger/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Fatalf("failed application run: %s", err)
		}
	}()

	gracefulShutdown(cancel)
}

func gracefulShutdown(cancel context.CancelFunc) {
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	var timeout = 5 * time.Second
	log.Printf("after %v seconds, the server will stop", timeout.Seconds())
	cancel()
	time.Sleep(timeout)
	os.Exit(0)
}

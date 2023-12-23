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
	go gracefulShutdown(cancel)
	runApplication(ctx)
}

func runApplication(ctx context.Context) {
	err := app.Run(ctx)
	if err != nil {
		log.Fatalf("application: %s", err)
	}
}

func gracefulShutdown(cancel context.CancelFunc) {
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	log.Printf("%s signal was received", <-quit)
	var timeout = 2 * time.Second
	log.Printf("after %v seconds, the program will force exit\n", timeout.Seconds())
	cancel()
	time.Sleep(timeout)
	os.Exit(0)
}

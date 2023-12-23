package main

import (
	"chatmerger/internal/app"
	"chatmerger/internal/config"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	cfg := initConfig()
	ctx, cancel := context.WithCancel(context.Background())
	go gracefulShutdown(cancel)
	runApplication(ctx, cfg)
}

func runApplication(ctx context.Context, cfg *config.Config) {
	err := app.Run(ctx, cfg)
	if err != nil {
		log.Fatalf("application: %s", err)
	}
}

func initConfig() *config.Config {
	cfgFs := config.InitFlagSet()

	cfg, err := cfgFs.Parse(os.Args[1:])
	if err != nil {
		log.Printf("config FlagSet initialization: %s", err)
		if errors.Is(err, config.WrongArgumentError) {
			cfgFs.FlagSetUsage()
		}
		os.Exit(2)
	}
	return cfg
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

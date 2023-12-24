package main

import (
	"chatmerger/internal/app"
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/common/vals"
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
	log.Println(msgs.ServerStarting)
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	cfg := initConfig()
	log.Println(msgs.ConfigInitialized)
	ctx, cancel := context.WithCancel(context.Background())
	go runApplication(ctx, cfg)
	gracefulShutdown(cancel)
}

func runApplication(ctx context.Context, cfg *config.Config) {
	log.Println(msgs.ApplicationStart)
	err := app.Run(ctx, cfg)
	if err != nil {
		log.Fatalf("application: %s", err)
	}
	os.Exit(0)
}

func initConfig() *config.Config {
	cfgFs := config.InitFlagSet()

	cfg, err := cfgFs.Parse(os.Args[1:])
	if err != nil {
		log.Printf("config FlagSet initialization: %s", err)
		if errors.Is(err, config.WrongArgumentError) {
			cfgFs.Usage()
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
	log.Println(msgs.ProgramWillForceExit)
	cancel()
	time.Sleep(vals.GracefulShutdownTimeout)
	os.Exit(0)
}

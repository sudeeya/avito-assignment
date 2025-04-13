package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/sudeeya/avito-assignment/internal/app"
	"github.com/sudeeya/avito-assignment/internal/config"
	"github.com/sudeeya/avito-assignment/internal/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("creating config: %v", err)
	}

	logger.SetGlobalLogger(cfg.LogConfig)

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		zap.S().Fatalf("creating app: %v", err)
	}

	a.Run(ctx)
}

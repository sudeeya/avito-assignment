package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sudeeya/avito-assignment/internal/config"
	"github.com/sudeeya/avito-assignment/internal/httpserver"
	"github.com/sudeeya/avito-assignment/internal/repository/postgres"
	"github.com/sudeeya/avito-assignment/internal/service"
	"go.uber.org/zap"
)

const (
	_shutdownTimeout = 5 * time.Second
)

type App struct {
	server *http.Server
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	repo, err := postgres.NewPostgres(ctx, cfg.DBConfig)
	if err != nil {
		return nil, fmt.Errorf("creating repository: %w", err)
	}

	services, err := service.NewService(cfg.ServerConfig, repo)
	if err != nil {
		return nil, fmt.Errorf("creating services: %w", err)
	}

	server := httpserver.NewServer(cfg.ServerConfig, services)

	return &App{
		server: server,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			zap.S().Errorf("Listen and serve: %v", err)
		}

		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		zap.L().Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), _shutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			zap.S().Errorf("Shutdown: %v", err)
		}
	}
}

package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/sudeeya/avito-assignment/internal/config"
	grpc_v1 "github.com/sudeeya/avito-assignment/internal/controller/grpc/v1"
	http_v1 "github.com/sudeeya/avito-assignment/internal/controller/http/v1"
	"github.com/sudeeya/avito-assignment/internal/grpcserver"
	"github.com/sudeeya/avito-assignment/internal/httpserver"
	"github.com/sudeeya/avito-assignment/internal/repository/postgres"
	"github.com/sudeeya/avito-assignment/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	_shutdownTimeout = 5 * time.Second
)

type App struct {
	cfg        *config.Config
	httpServer *http.Server
	grpcServer *grpc.Server
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

	router := http_v1.NewRouter(services)
	httpServer := httpserver.NewServer(cfg.ServerConfig, router)

	pvzServiceServer := grpc_v1.NewPVZServiceServerImplementation(services)
	grpcServer := grpcserver.NewServer(pvzServiceServer)

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	var (
		httpDone = make(chan struct{})
		grpcDone = make(chan struct{})
	)

	go func() {
		defer close(httpDone)

		zap.L().Info("Server is serving HTTP...")
		if err := a.httpServer.ListenAndServe(); err != nil {
			zap.S().Errorf("Serving HTTP: %v", err)
		}
	}()

	go func() {
		defer close(grpcDone)

		listener, err := net.Listen("tcp", ":"+strconv.Itoa(a.cfg.ServerConfig.ServerGRPCPort))
		if err != nil {
			zap.S().Errorf("Announcing gRPC: %v", err)
			return
		}

		zap.L().Info("Server is serving gRPC...")
		if err := a.grpcServer.Serve(listener); err != nil {
			zap.S().Errorf("Serving gRPC: %v", err)
		}
	}()

	select {
	case <-httpDone:
		a.Shutdown(ctx)
	case <-grpcDone:
		a.Shutdown(ctx)
	case <-ctx.Done():
		a.Shutdown(ctx)
	}
}

func (a *App) Shutdown(ctx context.Context) {
	zap.L().Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(ctx, _shutdownTimeout)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		zap.S().Errorf("HTTP Server shutdown: %v", err)
	}

	a.grpcServer.GracefulStop()
}

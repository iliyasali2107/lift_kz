package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"mado/internal/config"
	httphandler "mado/internal/controller/http/handler"
	"mado/internal/core"
	"mado/internal/repository/psql"
	"mado/pkg/database/postgres"
	"mado/pkg/httpserver"
)

// App is a application interface.
type App struct {
	logger     *zap.Logger
	db         *postgres.Postgres
	httpServer httpserver.Server
}

// New creates a new App.
func New(ctx context.Context, logger *zap.Logger) (App, error) {
	cfg := config.Get()

	postgresInstance, err := postgres.New(
		ctx,
		postgres.NewConnectionConfig(
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.DBName,
			cfg.Postgres.User,
			cfg.Postgres.Password,
			cfg.Postgres.SSLMode,
		),
	)
	if err != nil {
		return App{}, fmt.Errorf("can not connect to postgres: %w", err)
	}

	repositories := psql.NewRepositories(postgresInstance, logger)
	services := core.NewServices(repositories, logger)
	

	router := httphandler.NewRouter(httphandler.Deps{
		Logger:   logger,
		Services: services,
	})

	return App{
		logger: logger,
		db:     postgresInstance,
		httpServer: httpserver.New(
			router,
			httpserver.WithHost(cfg.HTTP.Host),
			httpserver.WithPort(cfg.HTTP.Port),
			httpserver.WithMaxHeaderBytes(cfg.HTTP.MaxHeaderBytes),
			httpserver.WithReadTimeout(cfg.HTTP.ReadTimeout),
			httpserver.WithWriteTimeout(cfg.HTTP.WriteTimeout),
		),
	}, nil
}

// Run runs the application.
func (a App) Run(ctx context.Context) error {
	eChan := make(chan error)
	interrupt := make(chan os.Signal, 1)

	a.logger.Info("Http server is starting")

	go func() {
		if err := a.httpServer.Start(); err != nil {
			eChan <- fmt.Errorf("failed to listen and serve: %w", err)
		}
	}()

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-eChan:
		return fmt.Errorf("conduit started failed: %w", err)
	case <-interrupt:
	}

	const httpShutdownTimeout = 5 * time.Second
	if err := a.httpServer.Stop(ctx, httpShutdownTimeout); err != nil {
		return fmt.Errorf("failed to stop http server: %w", err)
	}

	return nil
}

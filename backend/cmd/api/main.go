package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"mado/internal/app/api"
	"mado/internal/config"
	"mado/pkg/logger"
)

var (
	version   string
	buildDate string
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg := config.Get()
	l := logger.New(os.Stdout, logger.WithLevel(cfg.Logger.Level))
	l.Info("petition_service", zap.String("version", version), zap.String("build_date", buildDate))

	app, err := api.New(ctx, l)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}

package app

import (
	"companies/internal/config"
	"companies/internal/ctx"
	"context"
	"net/http"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	logger := InitializeLogger(cfg)
	ctx := ctx.WithLogger(context.Background(), logger)
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	mongo := InitializeMongo(ctx, cfg, logger)
	handler := InitializeRouter(cfg, logger, mongo)
	logger.Info("Starting HTTP server and serve port " + cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, handler)

	if err != nil {
		logger.Errorw("error starting HTTP listener", "error", err)
		return
	}
}

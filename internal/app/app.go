package app

import (
	"companies/internal/company_ctx"
	"companies/internal/config"
	"context"
	"net/http"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	logger := InitializeLogger(cfg)
	ctx := company_ctx.WithLogger(context.Background(), logger)
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	mongo := InitializeMongo(ctx, cfg, logger)
	handler := InitializeRouter(cfg, logger, mongo)

	if mongo != nil {
		defer func() {
			if err := mongo.Disconnect(context.Background()); err != nil {
				panic(err)
			}
		}()
		logger.Debug("mongodb initialized")
	} else {
		logger.Error("mongodb init failed")
	}

	logger.Info("Starting HTTP server and serve port " + cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, handler)

	if err != nil {
		logger.Errorw("error starting HTTP listener", "error", err)
		return
	}

}

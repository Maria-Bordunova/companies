//go:build wireinject
// +build wireinject

package app

import (
	"companies/internal/config"
	"companies/internal/domain/interfaces"
	"companies/internal/infra/http/auth"
	v1 "companies/internal/infra/http/v1"
	controller "companies/internal/infra/http/v1/controller"
	"companies/internal/infra/http/v1/middleware"
	"companies/pkg/logger"
	"context"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"

	repoMongo "companies/internal/infra/db/mongo"

	"net/http"
)

func InitializeLogger(cfg *config.Config) *logger.Logger {
	qloggerConfig := &cfg.Log
	qLogger := logger.New(qloggerConfig)
	return qLogger
}

func InitializeMongo(ctx context.Context, cfg *config.Config, logger *logger.Logger) *mongo.Client {
	panic(wire.Build(
		repoMongo.NewClient,
		wire.FieldsOf(&cfg, "Mongo"),
	))
}

func InitializeRouter(cfg *config.Config, logger *logger.Logger, mongo *mongo.Client) http.Handler {
	panic(wire.Build(
		// http
		middleware.NewRegistry,
		controller.NewController,
		v1.NewRouter,

		// auth
		wire.Bind(new(interfaces.UserAuthorizer), new(*auth.JwtValidator)),
		auth.NewJwtValidator,
		wire.FieldsOf(&cfg, "Auth"),

		// mongo
		wire.Bind(new(interfaces.CompaniesRepo), new(*repoMongo.CompaniesRepoMongo)),
		repoMongo.NewCompaniesRepo,

		// kafka

	))

	return nil
}

package config

import (
	"companies/internal/infra/db/mongo"
	"companies/pkg/logger"
)

type Auth struct {
	JwtKey string `env:"JWT_SECRET_KEY" default:"jwt_secret_key"`
}

type Config struct {
	Port  string `env:"APP_PORT" default:"80"`
	Mongo mongo.Config
	Auth  Auth
	Log   logger.Config
}

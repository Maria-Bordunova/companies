package config

import (
	"companies/internal/infra/db/mongo"
	"companies/pkg/logger"
)

type Auth struct {
	JwtKey string `env:"JWT_SECRET_KEY" default:"jwt_secret_key"`
}

type Kafka struct {
	Host  string `env:"KAFKA_HOST" default:"kafka:29092"`
	Topic string `env:"KAFKA_TOPIC" default:"company_changed"`
}

type Config struct {
	Port  string `env:"APP_PORT" default:"80"`
	Mongo mongo.Config
	Auth  Auth
	Log   logger.Config
	Kafka Kafka
}

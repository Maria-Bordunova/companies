package mongo

import (
	"companies/pkg/logger"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"time"
)

type Config struct {
	Dsn string `env:"MONGO_DSN"`
}

func NewClient(ctx context.Context, config Config, logger *logger.Logger) *mongo.Client {
	return connect(ctx, config, logger)
}

func connect(ctx context.Context, config Config, logger *logger.Logger) *mongo.Client {
	if config.Dsn == "" {
		logger.Errorf("MongoDB error: config: empty DSN")
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Dsn))
	if err != nil {
		logger.Errorf("MongoDB error: connecting: %v", err)
		return nil
	}

	parsedDSN, err := url.ParseRequestURI(config.Dsn)
	if err != nil {
		logger.Infof("Connecting to MongoDB: %s", config.Dsn)
	} else {
		logger.Infof("Connecting to MongoDB: %s", parsedDSN.Scheme+"://"+parsedDSN.Host)
	}

	// block until connection confirmed
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Errorf("MongoDB error: checking connect: %v", err)
		return nil
	}
	return client
}

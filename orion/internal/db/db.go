package db

import (
	"context"
	"github.com/andibalo/ramein/orion/ent"

	"github.com/andibalo/ramein/orion/internal/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func InitDB(cfg config.Config) *ent.Client {
	connStr := cfg.DBConnString()

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		cfg.Logger().Error("Failed to connect to db", zap.Error(err))
		panic("Failed to connect to db")
	}

	cfg.Logger().Info("Connected to database")

	if err := client.Schema.Create(context.Background()); err != nil {
		cfg.Logger().Error("failed creating schema resources", zap.Error(err))
		panic("failed creating schema resources")
	}

	return client
}

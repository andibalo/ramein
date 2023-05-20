package db

import (
	"context"
	"github.com/andibalo/ramein/phoenix/internal/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitDB(ctx context.Context, cfg config.Config) neo4j.DriverWithContext {
	uri := cfg.StorageConfig().Uri
	auth := neo4j.BasicAuth(cfg.StorageConfig().User, cfg.StorageConfig().Password, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		cfg.Logger().Error("Failed to init neo4j driver")
		panic(err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		cfg.Logger().Error("Failed to connect to neo4j db")
		panic(err)
	}

	cfg.Logger().Info("Connected to neo4j db")

	return driver
}

package db

import (
	"database/sql"
	"github.com/andibalo/ramein/core/internal/config"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func InitDB(cfg config.Config) *bun.DB {
	connStr := cfg.DBConnString()
	// open database

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))

	pgdb := bun.NewDB(db, pgdialect.New())

	if cfg.AppEnv() == "DEV" {
		pgdb.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	cfg.Logger().Info("Connected to database")

	return pgdb
}

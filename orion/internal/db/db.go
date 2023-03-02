package db

import (
	"entgo.io/ent/entc/integration/ent"
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
	defer client.Close()

	cfg.Logger().Info("Connected to database")

	//if err := client.Schema.Create(context.Background()); err != nil {
	//	log.Fatalf("failed creating schema resources: %v", err)
	//}

	return client
}

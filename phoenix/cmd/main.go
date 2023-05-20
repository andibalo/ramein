package main

import (
	"context"
	"fmt"
	"github.com/andibalo/ramein/phoenix"
	"github.com/andibalo/ramein/phoenix/internal/db"

	"github.com/andibalo/ramein/phoenix/internal/config"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	cfg := config.InitConfig()

	ctx := context.Background()

	driver := db.InitDB(ctx, cfg)

	defer driver.Close(ctx)

	server := phoenix.NewServer(ctx, cfg, driver)

	err = server.Start(cfg.AppAddress())

	if err != nil {
		cfg.Logger().Fatal("Port already used")
	}
}

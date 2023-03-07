package main

import (
	"fmt"
	"github.com/andibalo/ramein/orion"
	"github.com/andibalo/ramein/orion/internal/config"
	"github.com/andibalo/ramein/orion/internal/db"
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

	database := db.InitDB(cfg)

	defer database.Close()

	server := orion.NewServer(cfg, database)

	err = server.Start(cfg.AppAddress())

	if err != nil {
		cfg.Logger().Fatal("Port already used")
	}
}

package main

import (
	"fmt"
	"github.com/andibalo/ramein/phoenix"

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

	server := phoenix.NewServer(cfg)

	err = server.Start(cfg.AppAddress())

	if err != nil {
		cfg.Logger().Fatal("Port already used")
	}
}

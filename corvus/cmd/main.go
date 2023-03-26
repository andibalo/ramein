package main

import (
	"fmt"
	"github.com/andibalo/ramein/corvus"
	"github.com/andibalo/ramein/corvus/internal/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	srv := corvus.NewGRPCServer(cfg)

	cfg.Logger().Info(fmt.Sprintf("starting grpc server at port %v", cfg.AppAddress()))
	if err = srv.Start(); err != nil {
		cfg.Logger().Error(fmt.Sprintf("failed to start server"), zap.Error(err))
		panic("failed to start server")
	}
}

package config

import (
	"fmt"
	"github.com/andibalo/ramein/corvus/internal/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	AppPort            = ":8082"
	EnvDevEnvironment  = "DEV"
	EnvProdEnvironment = "PROD"
)

type Config interface {
	Logger() *zap.Logger

	AppEnv() string
	AppAddress() string
}

type AppConfig struct {
	logger *zap.Logger
	App    app
}

type app struct {
	AppEnv      string
	AppVersion  string
	Name        string
	Description string
	AppUrl      string
	AppID       string
}

func InitConfig() *AppConfig {
	viper.SetConfigType("env")
	viper.SetConfigName(".env") // name of Config file (without extension)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return &AppConfig{}
	}

	l := logger.InitLogger()

	return &AppConfig{
		logger: l,
		App: app{
			AppEnv:      viper.GetString("APP_ENV"),
			AppVersion:  viper.GetString("APP_VERSION"),
			Name:        "corvus",
			Description: "file manager service",
			AppUrl:      viper.GetString("APP_URL"),
			AppID:       viper.GetString("APP_ID"),
		},
	}
}

func getRequiredString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

func (c *AppConfig) Logger() *zap.Logger {
	return c.logger
}

func (c *AppConfig) AppEnv() string {
	return c.App.AppEnv
}

func (c *AppConfig) AppAddress() string {
	return c.App.AppUrl + AppPort
}

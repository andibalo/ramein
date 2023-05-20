package config

import (
	"fmt"
	"github.com/andibalo/ramein/phoenix/internal/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	AppAddress         = ":8085"
	EnvDevEnvironment  = "DEV"
	EnvProdEnvironment = "PROD"
)

type Config interface {
	Logger() *zap.Logger
	StorageConfig() db

	AppEnv() string
	AppAddress() string
	DbUserName() string
}

type AppConfig struct {
	logger *zap.Logger
	App    app
	Db     db
}

type app struct {
	AppEnv      string
	AppVersion  string
	Name        string
	Description string
	AppUrl      string
	AppID       string
}

type db struct {
	Uri      string
	User     string
	Password string
	Name     string
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
			Name:        "phoenix",
			Description: "friends management service",
			AppUrl:      viper.GetString("APP_URL"),
			AppID:       viper.GetString("APP_ID"),
		},
		Db: db{
			Uri:      getRequiredString("NEO4J_DB_URI"),
			User:     viper.GetString("NEO4J_DB_USER"),
			Password: viper.GetString("NEO4J_DB_PASSWORD"),
			Name:     viper.GetString("NEO4J_DB_NAME"),
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

func (c *AppConfig) StorageConfig() db {
	return c.Db
}

func (c *AppConfig) DbUserName() string {
	return c.Db.User
}

func (c *AppConfig) AppEnv() string {
	return c.App.AppEnv
}

func (c *AppConfig) AppAddress() string {
	return AppAddress
}

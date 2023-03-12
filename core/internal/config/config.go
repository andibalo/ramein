package config

import (
	"fmt"
	"github.com/andibalo/ramein/core/internal/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	AppName            = "CORE_SERVICE"
	AppAddress         = ":8000"
	EnvDevEnvironment  = "DEV"
	EnvProdEnvironment = "PROD"
)

type Config interface {
	Logger() *zap.Logger
	StorageConfig() db

	AppEnv() string
	AppAddress() string

	DBConnString() string

	RabbitMQURL() string
	RabbitMQChannel() string

	UserSecretCodeExpiryMins() int
}

type AppConfig struct {
	logger *zap.Logger
	App    app
	Db     db
	Rmq    rmq
}

type app struct {
	AppEnv                   string
	AppVersion               string
	Name                     string
	Description              string
	AppUrl                   string
	AppID                    string
	UserSecretCodeExpiryMins int
}

type db struct {
	DSN      string
	User     string
	Password string
	Name     string
	Host     string
	Port     int
	MaxPool  int
}

type rmq struct {
	Channel string
	URL     string
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
			AppEnv:                   viper.GetString("APP_ENV"),
			AppVersion:               viper.GetString("APP_VERSION"),
			Name:                     "core",
			Description:              "Stores and maintains user information",
			AppUrl:                   viper.GetString("APP_URL"),
			AppID:                    viper.GetString("APP_ID"),
			UserSecretCodeExpiryMins: viper.GetInt("USER_SECRET_CODE_EXPIRY_MINS"),
		},
		Db: db{
			DSN:      getRequiredString("DB_DSN"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			MaxPool:  viper.GetInt("DB_MAX_POOLING_CONNECTION"),
		},
		Rmq: rmq{
			Channel: viper.GetString("RABBITMQ_CHANNEL"),
			URL:     viper.GetString("RABBITMQ_URL"),
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

func (c *AppConfig) AppEnv() string {
	return c.App.AppEnv
}

func (c *AppConfig) AppAddress() string {
	return AppAddress
}

func (c *AppConfig) DBConnString() string {
	return c.StorageConfig().DSN
}

func (c *AppConfig) UserSecretCodeExpiryMins() int {
	return c.App.UserSecretCodeExpiryMins
}

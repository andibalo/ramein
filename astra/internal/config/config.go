package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

const (
	AppPort            = ":8086"
	EnvDevEnvironment  = "DEV"
	EnvProdEnvironment = "PROD"
)

type Config interface {
	Logger() *zap.Logger

	AppEnv() string
	AppAddress() string

	DBKeyspace() string
	DBHosts() []string

	KafkaHosts() []string

	KafkaPendingMessagesTopic() string
}

type AppConfig struct {
	logger      *zap.Logger
	App         app
	Db          db
	Kafka       kafka
	KafkaTopics kafkaTopics
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
	hosts    string
	keyspace string
}

type kafka struct {
	hosts string
}

type kafkaTopics struct {
	pendingMessagesTopic string
}

func InitConfig(logger *zap.Logger) *AppConfig {
	return &AppConfig{
		logger: logger,
		App: app{
			AppEnv:      viper.GetString("APP_ENV"),
			AppVersion:  viper.GetString("APP_VERSION"),
			Name:        "astra",
			Description: "chat service",
			AppUrl:      viper.GetString("APP_URL"),
			AppID:       viper.GetString("APP_ID"),
		},
		Db: db{
			hosts:    viper.GetString("DB_HOSTS"),
			keyspace: viper.GetString("DB_KEYSPACE"),
		},
		Kafka: kafka{
			hosts: viper.GetString("KAFKA_HOSTS"),
		},
		KafkaTopics: kafkaTopics{
			pendingMessagesTopic: viper.GetString("KAFKA_PENDING_MESSAGES_TOPIC"),
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

func (c *AppConfig) DBHosts() []string {

	dbHosts := strings.Split(c.Db.hosts, ",")

	return dbHosts
}

func (c *AppConfig) DBKeyspace() string {
	return c.Db.keyspace
}

func (c *AppConfig) KafkaHosts() []string {

	kafkaHosts := strings.Split(c.Kafka.hosts, ",")

	return kafkaHosts
}

func (c *AppConfig) KafkaPendingMessagesTopic() string {
	return c.KafkaTopics.pendingMessagesTopic
}

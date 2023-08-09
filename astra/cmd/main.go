package main

import (
	"context"
	"fmt"
	"github.com/andibalo/ramein/astra/internal/config"
	"github.com/andibalo/ramein/astra/internal/db"
	"github.com/andibalo/ramein/astra/internal/logger"
	"github.com/andibalo/ramein/astra/internal/model"
	"github.com/andibalo/ramein/astra/internal/redis"
	"github.com/andibalo/ramein/astra/internal/repository"
	"github.com/andibalo/ramein/commons/kafka"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	l := logger.InitLogger()

	cfg := config.InitConfig(l)

	session, err := db.InitDB(cfg)
	if err != nil {
		panic(err)
	}

	err = db.InitKeyspaceAndTables(cfg, session)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	cache := redis.NewRedisDriver(cfg)

	_, err = cache.Connect(context.Background())
	if err != nil {
		l.Error("error init redis cache", zap.Error(err))
	}

	pendingMessagesTopicSyncProducer, err := kafka.NewSyncProducer(
		cfg.KafkaHosts(),
		cfg.KafkaPendingMessagesTopic(),
		kafka.WithLogger(l),
	)

	if err != nil {
		l.Error("error init pending_messages kafka producer", zap.Error(err))
	}

	defer pendingMessagesTopicSyncProducer.Close()

	messageRepo := repository.NewMessageRepository(session, cfg.Logger())

	r := gin.Default()

	r.GET("/test-insert", func(c *gin.Context) {

		conversationID, _ := gocql.RandomUUID()
		messageID, _ := gocql.RandomUUID()

		m := model.Message{
			ConversationID:    conversationID,
			MessageID:         messageID,
			ConversationName:  "test",
			FromUserID:        "user-id",
			FromUserNumber:    "0929024",
			FromUserFirstName: "andi",
			FromUserLastName:  "balo",
			FromUserEmail:     "andialo214@gmail.com",
			TextContent:       "hello",
			SentAt:            time.Now(),
			CreatedBy:         "",
			CreatedAt:         time.Now(),
			UpdatedBy:         "",
			UpdatedAt:         nil,
			DeletedBy:         "",
			DeletedAt:         nil,
		}

		err = messageRepo.SaveMessage(m)

		if err != nil {
			c.String(500, "fail")
		}

		c.String(200, "success")
	})

	_ = r.Run() // listen and serve on 0.0.0.0:8080
}

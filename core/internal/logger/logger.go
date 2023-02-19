package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

func InitLogger() *zap.Logger {
	var logger *zap.Logger

	logger, err := zap.NewProduction()

	if err != nil {
		panic(fmt.Sprintf("logger initialization failed %v", err))
	}

	if os.Getenv("APP_ENV") == "DEV" {
		logger, err = zap.NewDevelopment()

		if err != nil {
			panic(fmt.Sprintf("logger initialization failed %v", err))
		}
	}

	logger.Info("logger started")

	defer logger.Sync()

	return logger
}

package core

import (
	"github.com/andibalo/ramein/commons/rabbitmq"
	"github.com/andibalo/ramein/core/internal/api"
	v1 "github.com/andibalo/ramein/core/internal/api/v1"
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/andibalo/ramein/core/internal/httpresp"
	"github.com/andibalo/ramein/core/internal/pubsub"
	"github.com/andibalo/ramein/core/internal/repository"
	"github.com/andibalo/ramein/core/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/uptrace/bun"
	"time"
)

const idleTimeout = 5 * time.Second

func NewServer(cfg config.Config, db *bun.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  idleTimeout,
		ErrorHandler: httpresp.HttpRespError,
	})

	app.Use(recover.New())

	rmq := rabbitmq.NewRabitmq(rabbitmq.RabitmqConfiguration{
		URL:    cfg.RabbitMQURL(),
		Enable: true,
	})

	pb := pubsub.NewPubSub(cfg, rmq)

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(cfg, userRepo, pb)

	userController := v1.NewUserController(cfg, userService)

	registerHandlers(app, &api.HealthCheck{}, userController)

	return app
}

func registerHandlers(e *fiber.App, handlers ...api.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(e)
	}
}

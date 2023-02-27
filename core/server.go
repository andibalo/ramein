package core

import (
	"github.com/andibalo/ramein/core/internal/api"
	v1 "github.com/andibalo/ramein/core/internal/api/v1"
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/andibalo/ramein/core/internal/httpresp"
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
	userController := v1.NewUserController(cfg)

	registerHandlers(app, &api.HealthCheck{}, userController)

	return app
}

func registerHandlers(e *fiber.App, handlers ...api.Handler) {
	for _, handler := range handlers {
		handler.AddRoutes(e)
	}
}

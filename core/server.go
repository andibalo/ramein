package core

import (
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"time"
)

const idleTimeout = 5 * time.Second

func NewServer(cfg config.Config, db *bun.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	return app
}

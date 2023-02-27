package api

import "github.com/gofiber/fiber/v2"

type Handler interface {
	AddRoutes(f *fiber.App)
}

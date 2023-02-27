package api

import (
	"github.com/gofiber/fiber/v2"
)

const (
	healthCheckPath = "/health"
)

// HealthCheck is a standard, simple health check
type HealthCheck struct{}

// AddRoutes adds the routers for this API to the provided router (or subrouter)
func (h *HealthCheck) AddRoutes(f *fiber.App) {
	f.Get(healthCheckPath, h.handler)
}

func (h *HealthCheck) handler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("OK")
}

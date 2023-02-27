package v1

import (
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/andibalo/ramein/core/internal/constants"
	"github.com/gofiber/fiber/v2"
)

const (
	userBasePath = "/user"
)

type UserController struct {
	cfg config.Config
}

func NewUserController(cfg config.Config) *UserController {

	return &UserController{
		cfg: cfg,
	}
}

func (h *UserController) AddRoutes(f *fiber.App) {
	r := f.Group(constants.V1BasePath + userBasePath)

	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.Login)
}

func (h *UserController) RegisterUser(c *fiber.Ctx) error {
	return c.SendString("Register User")
}

func (h *UserController) Login(c *fiber.Ctx) error {
	return c.SendString("Login User")
}

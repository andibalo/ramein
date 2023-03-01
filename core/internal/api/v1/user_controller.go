package v1

import (
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/andibalo/ramein/core/internal/constants"
	"github.com/andibalo/ramein/core/internal/httpresp"
	"github.com/andibalo/ramein/core/internal/request"
	"github.com/andibalo/ramein/core/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
)

const (
	userBasePath = "/user"
)

type UserController struct {
	cfg         config.Config
	userService service.UserService
}

func NewUserController(cfg config.Config, userService service.UserService) *UserController {

	return &UserController{
		cfg:         cfg,
		userService: userService,
	}
}

func (h *UserController) AddRoutes(f *fiber.App) {
	r := f.Group(constants.V1BasePath + userBasePath)

	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.Login)
}

func (h *UserController) RegisterUser(c *fiber.Ctx) error {
	req := &request.RegisterUserRequest{}

	if err := c.BodyParser(req); err != nil {
		h.cfg.Logger().Error("[RegisterUser] Failed to parse request body", zap.Error(err))
		return httpresp.HttpRespError(c, fiber.NewError(http.StatusBadRequest, "Failed to parse request body"))
	}

	err := h.userService.CreateUser(req)
	if err != nil {
		h.cfg.Logger().Error("[RegisterUser] Failed to parse request body", zap.Error(err))
		return httpresp.HttpRespError(c, err)
	}

	return httpresp.HttpRespSuccess(c, nil, nil)
}

func (h *UserController) Login(c *fiber.Ctx) error {
	return c.SendString("Login User")
}

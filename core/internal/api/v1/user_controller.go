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
	r := f.Group(constants.UserBasePathV1)

	r.Post(constants.UserRegisterPath, h.RegisterUser)
	r.Post(constants.UserLoginPath, h.Login)
	r.Get(constants.UserVerifyEmailPath, h.VerifyEmail)
}

func (h *UserController) RegisterUser(c *fiber.Ctx) error {
	req := &request.RegisterUserRequest{}

	if err := c.BodyParser(req); err != nil {
		h.cfg.Logger().Error("[RegisterUser] Failed to parse request body", zap.Error(err))
		return httpresp.HttpRespError(c, fiber.NewError(http.StatusBadRequest, "Failed to parse request body"))
	}

	err := h.userService.CreateUser(req)
	if err != nil {
		h.cfg.Logger().Error("[RegisterUser] Failed to create user", zap.Error(err))
		return httpresp.HttpRespError(c, err)
	}

	return httpresp.HttpRespSuccess(c, nil, nil)
}

func (h *UserController) Login(c *fiber.Ctx) error {
	req := &request.LoginRequest{}

	if err := c.BodyParser(req); err != nil {
		h.cfg.Logger().Error("[Login] Failed to parse request body", zap.Error(err))
		return httpresp.HttpRespError(c, fiber.NewError(http.StatusBadRequest, "Failed to parse request body"))
	}

	jwt, err := h.userService.Login(req)
	if err != nil {
		h.cfg.Logger().Error("[Login] Failed to login user", zap.Error(err))
		return httpresp.HttpRespError(c, err)
	}

	return httpresp.HttpRespSuccess(c, jwt, nil)
}

func (h *UserController) VerifyEmail(c *fiber.Ctx) error {
	secretCode := c.Query("secret_code")
	userVerifyId := c.Query("id")

	if secretCode == "" || userVerifyId == "" {
		h.cfg.Logger().Error("[VerifyEmail] Secret code and user verify id query param must exist")
		return httpresp.HttpRespError(c, fiber.NewError(http.StatusBadRequest, "Invalid verify email link"))
	}
	//
	//jwt, err := h.userService.Login(req)
	//if err != nil {
	//	h.cfg.Logger().Error("[Login] Failed to login user", zap.Error(err))
	//	return httpresp.HttpRespError(c, err)
	//}

	//return httpresp.HttpRespSuccess(c, jwt, nil)
	return nil
}

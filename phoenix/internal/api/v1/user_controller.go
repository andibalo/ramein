package v1

import (
	"github.com/andibalo/ramein/phoenix/internal/apperr"
	"github.com/andibalo/ramein/phoenix/internal/config"
	"github.com/andibalo/ramein/phoenix/internal/constants"
	"github.com/andibalo/ramein/phoenix/internal/httpresp"
	"github.com/andibalo/ramein/phoenix/internal/request"
	"github.com/andibalo/ramein/phoenix/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (h *UserController) AddRoutes(r *gin.Engine) {
	uc := r.Group(constants.V1BasePath + userBasePath)

	uc.GET("/", h.GetUsersList)
	uc.POST("/friend/request", h.SendFriendRequest)
	uc.POST("/friend/request/accept", h.AcceptFriendRequest)
	uc.GET("/friend/list/:user_id", h.GetFriendsList)
}

func (h *UserController) GetUsersList(c *gin.Context) {

	var req request.GetUsersListReq

	if err := c.BindQuery(&req); err != nil {
		httpresp.HttpRespError(c, apperr.ErrBadRequest)
		return
	}

	users, err := h.userService.GetUsersList(req)

	if err != nil {
		h.cfg.Logger().Error("[GetUsersList] Error fetching users list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, users, nil)
}

func (h *UserController) SendFriendRequest(c *gin.Context) {

	var req request.SendFriendRequestReq

	if err := c.BindJSON(&req); err != nil {
		httpresp.HttpRespError(c, apperr.ErrBadRequest)
		return
	}

	err := h.userService.SendFriendRequest(req)

	if err != nil {
		h.cfg.Logger().Error("[SendFriendRequest] Error sending friend request", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
}

func (h *UserController) AcceptFriendRequest(c *gin.Context) {

	var req request.AcceptFriendRequestReq

	if err := c.BindJSON(&req); err != nil {
		httpresp.HttpRespError(c, apperr.ErrBadRequest)
		return
	}

	err := h.userService.AcceptFriendRequest(req)

	if err != nil {
		h.cfg.Logger().Error("[AcceptFriendRequest] Error accepting friend request", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, nil, nil)
}

func (h *UserController) GetFriendsList(c *gin.Context) {

	var req request.GetFriendsListReq

	if err := c.BindQuery(&req); err != nil {
		httpresp.HttpRespError(c, apperr.ErrBadRequest)
		return
	}

	users, err := h.userService.GetFriendsList(c.Param("user_id"), req)

	if err != nil {
		h.cfg.Logger().Error("[GetFriendsList] Error fetching friends list", zap.Error(err))
		httpresp.HttpRespError(c, err)
		return
	}

	httpresp.HttpRespSuccess(c, users, nil)
}

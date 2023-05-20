package service

import (
	"github.com/andibalo/ramein/phoenix/internal/model"
	"github.com/andibalo/ramein/phoenix/internal/request"
)

type UserService interface {
	GetUsersList(req request.GetUsersListReq) ([]model.User, error)
}

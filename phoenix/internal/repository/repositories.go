package repository

import (
	"github.com/andibalo/ramein/phoenix/internal/model"
	"github.com/andibalo/ramein/phoenix/internal/request"
)

type UserRepository interface {
	FetchUsers(req request.GetUsersListReq) ([]model.User, error)
}

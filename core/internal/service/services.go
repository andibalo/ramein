package service

import "github.com/andibalo/ramein/core/internal/request"

type UserService interface {
	CreateUser(data *request.RegisterUserRequest) error
}
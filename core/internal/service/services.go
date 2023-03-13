package service

import "github.com/andibalo/ramein/core/internal/request"

type UserService interface {
	CreateUser(data *request.RegisterUserRequest) error
	Login(data *request.LoginRequest) (string, error)
	VerifyEmail(secretCode string, id string) error
}

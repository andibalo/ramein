package repository

import "github.com/andibalo/ramein/core/internal/model"

type UserRepository interface {
	Save(user model.User) error
}

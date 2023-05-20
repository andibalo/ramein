package service

import (
	"github.com/andibalo/ramein/phoenix/internal/config"
	"github.com/andibalo/ramein/phoenix/internal/model"
	"github.com/andibalo/ramein/phoenix/internal/repository"
	"github.com/andibalo/ramein/phoenix/internal/request"
	"go.uber.org/zap"
)

type userService struct {
	cfg      config.Config
	userRepo repository.UserRepository
}

func NewUserService(cfg config.Config, userRepo repository.UserRepository) *userService {
	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (s *userService) GetUsersList(req request.GetUsersListReq) ([]model.User, error) {

	users, err := s.userRepo.FetchUsers(req)
	if err != nil {
		s.cfg.Logger().Error("[GetUsersList] Error inserting template to db", zap.Error(err))
		return nil, err
	}

	return users, nil
}

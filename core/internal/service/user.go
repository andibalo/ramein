package service

import (
	"database/sql"
	"errors"
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/andibalo/ramein/core/internal/constants"
	"github.com/andibalo/ramein/core/internal/model"
	"github.com/andibalo/ramein/core/internal/pubsub"
	"github.com/andibalo/ramein/core/internal/repository"
	"github.com/andibalo/ramein/core/internal/request"
	"github.com/andibalo/ramein/core/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type userService struct {
	cfg      config.Config
	userRepo repository.UserRepository
	pb       pubsub.PubSub
}

func NewUserService(cfg config.Config, userRepo repository.UserRepository, pb pubsub.PubSub) *userService {

	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
		pb:       pb,
	}
}

func (s *userService) CreateUser(data *request.RegisterUserRequest) error {

	existingUser, err := s.userRepo.GetByEmail(data.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().Error("[CreateUser] Failed to get user by email", zap.Error(err))
		return err
	}

	if existingUser != nil {
		s.cfg.Logger().Error("[CreateUser] User already exists")
		return fiber.NewError(fiber.StatusBadRequest, "User already exists")
	}

	user, err := s.mapCreateUserReqToUserModel(data)
	if err != nil {
		s.cfg.Logger().Error("[CreateUser] Failed to map payload to user model", zap.Error(err))
		return err
	}

	err = s.userRepo.Save(user)
	if err != nil {
		s.cfg.Logger().Error("[CreateUser] Failed to insert user to database", zap.Error(err))
		return fiber.NewError(http.StatusInternalServerError, "Failed to insert user to database")
	}

	return nil
}

func (s *userService) mapCreateUserReqToUserModel(data *request.RegisterUserRequest) (*model.User, error) {

	hasedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		s.cfg.Logger().Error("[mapCreateUserReqToUserModel] Failed to hash password", zap.Error(err))

		return nil, fiber.NewError(http.StatusInternalServerError, "Failed to hash password")
	}

	id := constants.USER_ROLE_PREFIX + uuid.NewString()

	return &model.User{
		ID:              id,
		Email:           data.Email,
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		Phone:           data.Phone,
		Role:            data.Role,
		Password:        hasedPassword,
		IsSuperUser:     false,
		IsVerified:      false,
		IsEmailVerified: false,
		ProfileSummary:  data.ProfileSummary,
	}, nil
}

func (s *userService) Login(data *request.LoginRequest) (string, error) {

	existingUser, err := s.userRepo.GetByEmail(data.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().Error("[Login] Invalid email/password", zap.Error(err))
			return "", fiber.NewError(http.StatusBadRequest, "Invalid email/password")
		}

		s.cfg.Logger().Error("[Login] Failed to get user by email", zap.Error(err))
		return "", err
	}

	isMatch := util.CheckPasswordHash(data.Password, existingUser.Password)
	if !isMatch {
		s.cfg.Logger().Error("[Login] Invalid password for user", zap.String("email", data.Email))
		return "", fiber.NewError(http.StatusBadRequest, "Invalid email/password")
	}

	token, err := util.GenerateToken(existingUser)
	if err != nil {
		s.cfg.Logger().Error("[Login] Failed to generate JWT Token for user", zap.String("email", data.Email))
		return "", fiber.NewError(http.StatusInternalServerError, "Failed to generate JWT Token")
	}

	return token, nil
}

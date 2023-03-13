package service

import (
	"database/sql"
	"errors"
	"fmt"
	pubsubCommons "github.com/andibalo/ramein/commons/pubsub"
	"github.com/andibalo/ramein/core/internal/config"
	"github.com/andibalo/ramein/core/internal/constants"
	"github.com/andibalo/ramein/core/internal/model"
	"github.com/andibalo/ramein/core/internal/pubsub"
	"github.com/andibalo/ramein/core/internal/repository"
	"github.com/andibalo/ramein/core/internal/request"
	"github.com/andibalo/ramein/core/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type userService struct {
	cfg      config.Config
	userRepo repository.UserRepository
	pb       pubsub.PubSub
	db       *bun.DB
}

func NewUserService(cfg config.Config, userRepo repository.UserRepository, pb pubsub.PubSub, db *bun.DB) *userService {

	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
		pb:       pb,
		db:       db,
	}
}

func (s *userService) CreateUser(data *request.RegisterUserRequest) error {

	existingUser, err := s.userRepo.GetByEmail(data.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.cfg.Logger().Error("[CreateUser] Failed to get user by email", zap.Error(err))
		return fiber.NewError(http.StatusInternalServerError, "Failed toget user by email")
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

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().Error("[CreateUser] Failed to begin transaction", zap.Error(err))
		return fiber.NewError(http.StatusInternalServerError, "Failed to begin transaction")
	}

	err = s.userRepo.SaveTx(user, tx)
	if err != nil {
		s.cfg.Logger().Error("[CreateUser] Failed to insert user to database", zap.Error(err))
		tx.Rollback()

		return fiber.NewError(http.StatusInternalServerError, "Failed to insert user to database")
	}

	userVerifyEmail := &model.UserVerifyEmail{
		ID:         uuid.NewString(),
		UserID:     user.ID,
		SecretCode: util.GenRandomString(10),
		Email:      user.Email,
		IsUsed:     false,
		ExpiredAt:  time.Now().Add(time.Minute * time.Duration(s.cfg.UserSecretCodeExpiryMins())),
		CreatedBy:  config.AppName,
	}

	err = s.userRepo.SaveUserVerifyEmailTx(userVerifyEmail, tx)
	if err != nil {
		s.cfg.Logger().Error("[CreateUser] Failed to insert user verify email to database", zap.Error(err))
		tx.Rollback()

		return fiber.NewError(http.StatusInternalServerError, "Failed to insert user verify email to database")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().Error("[CreateUser] Failed to commit transaction", zap.Error(err))
		return fiber.NewError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	verifyUrl := s.cfg.AppURL() + constants.UserVerifyEmailPath + fmt.Sprintf("?secret_code=%s&id=%s", userVerifyEmail.SecretCode, userVerifyEmail.ID)

	msg := pubsubCommons.CoreNewRegisteredUserPayload{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		VerifyURL: verifyUrl,
	}

	go func() {
		s.pb.PublishNewUserRegistered(msg)
	}()

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
		return "", fiber.NewError(http.StatusInternalServerError, "Failed to get user by email")
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

func (s *userService) VerifyEmail(secretCode string, id string) error {

	uve, err := s.userRepo.GetUserVerifyEmailByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.cfg.Logger().Error("[VerifyEmail] User verify email not found", zap.Error(err))

			return fiber.NewError(http.StatusBadRequest, "Not found")
		}

		s.cfg.Logger().Error("[VerifyEmail] Failed to get user verify email by id", zap.Error(err))

		return fiber.NewError(http.StatusInternalServerError, "Failed to get user verify email by id")
	}

	if uve.IsUsed {
		s.cfg.Logger().Error("[VerifyEmail] Email is already verified", zap.Error(err))

		return fiber.NewError(http.StatusBadRequest, "Email is already verified")
	}

	if uve.SecretCode != secretCode {
		s.cfg.Logger().Error("[VerifyEmail] Invalid secret code", zap.Error(err))

		return fiber.NewError(http.StatusBadRequest, "Invalid email link")
	}

	if uve.ExpiredAt.Before(time.Now()) {
		s.cfg.Logger().Error("[VerifyEmail] Verify email link is expired", zap.Error(err))

		return fiber.NewError(http.StatusBadRequest, "Verify email link is expired")
	}

	tx, err := s.db.Begin()
	if err != nil {
		s.cfg.Logger().Error("[CreateUser] Failed to begin transaction", zap.Error(err))
		return fiber.NewError(http.StatusInternalServerError, "Failed to begin transaction")
	}

	err = s.userRepo.SetUserVerifyEmailToUsedTx(id, tx)
	if err != nil {
		s.cfg.Logger().Error("[VerifyEmail] Failed to set user verify email to used", zap.Error(err))
		tx.Rollback()

		return fiber.NewError(http.StatusInternalServerError, "Failed to set user verify email to used")
	}

	err = s.userRepo.SetUserToEmailVerifiedTx(uve.UserID, tx)
	if err != nil {
		s.cfg.Logger().Error("[VerifyEmail] Failed to set user to email verified", zap.Error(err))
		tx.Rollback()

		return fiber.NewError(http.StatusInternalServerError, "Failed to set user to email verified")
	}

	err = tx.Commit()
	if err != nil {
		s.cfg.Logger().Error("[VerifyEmail] Failed to commit transaction", zap.Error(err))
		return fiber.NewError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return nil
}

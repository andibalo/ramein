package service

import (
	"errors"
	"github.com/andibalo/ramein/phoenix/internal/apperr"
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

func (s *userService) SendFriendRequest(req request.SendFriendRequestReq) error {

	_, err := s.userRepo.FetchUserByUserID(req.UserID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			s.cfg.Logger().Error("[SendFriendRequest] Requester user not found", zap.Error(err))
			return err
		}

		s.cfg.Logger().Error("[SendFriendRequest] Error fetching requester user", zap.Error(err))
		return err
	}

	_, err = s.userRepo.FetchUserByUserID(req.TargetUserID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			s.cfg.Logger().Error("[SendFriendRequest] Target user not found", zap.Error(err))
			return err
		}

		s.cfg.Logger().Error("[SendFriendRequest] Error fetching target user", zap.Error(err))
		return err
	}

	err = s.userRepo.SaveFriendRequestRelationship(req.UserID, req.TargetUserID)
	if err != nil {
		s.cfg.Logger().Error("[SendFriendRequest] Error creating friend request relationship", zap.Error(err))
		return err
	}

	return nil
}

func (s *userService) AcceptFriendRequest(req request.AcceptFriendRequestReq) error {

	isAlreadyFriends, err := s.userRepo.CheckIsFriendsWithRelationshipExist(req.UserID, req.TargetUserID)
	if err != nil {
		if !errors.Is(err, apperr.ErrNotFound) {
			s.cfg.Logger().Error("[AcceptFriendRequest] Error fetching IS_FRIENDS_WITH relationship", zap.Error(err))
			return err
		}
	}

	if isAlreadyFriends {
		s.cfg.Logger().Info("[AcceptFriendRequest] User is already friends with target user")
		return nil
	}

	err = s.userRepo.SaveIsFriendsWithRelationship(req.UserID, req.TargetUserID)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			s.cfg.Logger().Error("[AcceptFriendRequest] User has not sent friend request to target user", zap.Error(err))
			return errors.New("User has not sent friend request")
		}

		s.cfg.Logger().Error("[AcceptFriendRequest] Error creating is friends with relationship", zap.Error(err))
		return err
	}

	return nil
}

func (s *userService) GetFriendsList(userID string, req request.GetFriendsListReq) ([]model.User, error) {

	userFriends, err := s.userRepo.FetchFriendsListByUserID(userID, req)
	if err != nil {
		s.cfg.Logger().Error("[GetFriendsList] Error fetching user friends list", zap.Error(err))
		return nil, err
	}

	return userFriends, nil
}

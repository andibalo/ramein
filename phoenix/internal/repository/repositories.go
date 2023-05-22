package repository

import (
	"github.com/andibalo/ramein/phoenix/internal/model"
	"github.com/andibalo/ramein/phoenix/internal/request"
)

type UserRepository interface {
	FetchUsers(req request.GetUsersListReq) ([]model.User, error)
	FetchUserByUserID(coreUserID string) (model.User, error)
	FetchFriendsListByUserID(userID string, req request.GetFriendsListReq) ([]model.User, error)
	SaveFriendRequestRelationship(userID, targetUserID string) error
	SaveIsFriendsWithRelationship(userID, targetUserID string) error
	CheckIsFriendsWithRelationshipExist(userID, targetUserID string) (bool, error)
}

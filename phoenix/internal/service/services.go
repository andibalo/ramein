package service

import (
	"github.com/andibalo/ramein/phoenix/internal/httpresp"
	"github.com/andibalo/ramein/phoenix/internal/model"
	"github.com/andibalo/ramein/phoenix/internal/request"
)

type UserService interface {
	GetUsersList(req request.GetUsersListReq) ([]model.User, error)
	SendFriendRequest(req request.SendFriendRequestReq) error
	AcceptFriendRequest(req request.AcceptFriendRequestReq) error
	GetFriendsList(userID string, req request.GetFriendsListReq) ([]model.User, *httpresp.Pagination, error)
}

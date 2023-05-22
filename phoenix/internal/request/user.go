package request

type GetUsersListReq struct {
	Limit int64 `json:"limit" form:"limit"`
}

type SendFriendRequestReq struct {
	UserID       string `json:"user_id"`
	TargetUserID string `json:"target_user_id"`
}

type AcceptFriendRequestReq struct {
	UserID       string `json:"user_id"`
	TargetUserID string `json:"target_user_id"`
}

type GetFriendsListReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

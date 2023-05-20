package request

type GetUsersListReq struct {
	Limit int64 `json:"limit" form:"limit"`
}

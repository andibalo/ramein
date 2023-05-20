package model

type User struct {
	ID          int64  `json:"id"`
	CoreUserID  string `json:"core_user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

package request

type RegisterUserRequest struct {
	Email          string   `json:"email"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Phone          string   `json:"phone"`
	Role           string   `json:"role"`
	Password       string   `json:"password"`
	ProfileSummary *string  `json:"profile_summary"`
	Images         []string `json:"images"`
}

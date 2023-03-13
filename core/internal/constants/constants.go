package constants

const (
	V1BasePath = "/api/v1"

	//user routes
	UserBasePathV1      = V1BasePath + "/user"
	UserRegisterPath    = "/register"
	UserLoginPath       = "/login"
	UserVerifyEmailPath = "/verify/email"
)

// role
const (
	USER             = "USER"
	ADMIN            = "ADMIN"
	USER_ROLE_PREFIX = "USR-"
)

package constants

const (
	V1BasePath = "/api/v1"

	//user routes
	UserBasePathV1      = V1BasePath + "/user"
	UserRegisterPath    = UserBasePathV1 + "/register"
	UserLoginPath       = UserBasePathV1 + "/login"
	UserVerifyEmailPath = UserBasePathV1 + "/verify/email"
)

// role
const (
	USER             = "USER"
	ADMIN            = "ADMIN"
	USER_ROLE_PREFIX = "USR-"
)

package response

const (
	Success           Code = "PHX0000"
	ServerError       Code = "PHX0001"
	BadRequest        Code = "PHX0002"
	InvalidRequest    Code = "PHX0004"
	Failed            Code = "PHX0073"
	Pending           Code = "PHX0050"
	InvalidInputParam Code = "PHX0032"
	DuplicateUser     Code = "PHX0033"
	NotFound          Code = "PHX0034"

	Unauthorized   Code = "PHX0502"
	Forbidden      Code = "PHX0503"
	GatewayTimeout Code = "PHX0048"
)

type Code string

var codeMap = map[Code]string{
	Success:           "success",
	Failed:            "failed",
	Pending:           "pending",
	BadRequest:        "bad or invalid request",
	Unauthorized:      "Unauthorized Token",
	GatewayTimeout:    "Gateway Timeout",
	ServerError:       "Internal Server Error",
	InvalidInputParam: "Other invalid argument",
	DuplicateUser:     "duplicate user",
	NotFound:          "Not found",
}

func (c Code) GetStatus() string {
	switch c {
	case Success:
		return "SUCCESS"

	default:
		return "FAILED"
	}
}

func (c Code) GetMessage() string {
	return codeMap[c]
}

func (c Code) GetVersion() string {
	return "1"
}

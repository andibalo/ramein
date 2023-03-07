package response

const (
	Success           Code = "OR0000"
	ServerError       Code = "OR0001"
	BadRequest        Code = "OR0002"
	InvalidRequest    Code = "OR0004"
	Failed            Code = "OR0073"
	Pending           Code = "OR0050"
	InvalidInputParam Code = "OR0032"
	DuplicateUser     Code = "OR0033"
	NotFound          Code = "OR0034"

	Unauthorized   Code = "OR0502"
	Forbidden      Code = "OR0503"
	GatewayTimeout Code = "OR0048"
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

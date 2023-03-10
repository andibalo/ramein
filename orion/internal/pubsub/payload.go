package pubsub

type CoreNewRegisteredUserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	VerifyURL string `json:"verify_url"`
}

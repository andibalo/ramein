package request

type CreateTemplateReq struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template"`
}

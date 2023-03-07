package service

import "github.com/andibalo/ramein/orion/internal/request"

type TemplateService interface {
	CreateTemplate(data request.CreateTemplateReq) error
}

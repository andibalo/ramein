package service

import (
	"github.com/andibalo/ramein/orion/ent"
	"github.com/andibalo/ramein/orion/internal/request"
)

type TemplateService interface {
	CreateTemplate(data request.CreateTemplateReq) error
	GetTemplateByID(templateID string) (*ent.Template, error)
}

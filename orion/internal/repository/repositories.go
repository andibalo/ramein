package repository

import (
	"github.com/andibalo/ramein/orion/ent"
	"github.com/andibalo/ramein/orion/internal/request"
)

type TemplateRepository interface {
	Save(data request.CreateTemplateReq) error
	GetByID(templateID string) (*ent.Template, error)
	GetByTemplateName(templateName string) (*ent.Template, error)
}

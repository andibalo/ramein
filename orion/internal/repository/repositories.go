package repository

import "github.com/andibalo/ramein/orion/internal/request"

type TemplateRepository interface {
	Save(data request.CreateTemplateReq) error
}

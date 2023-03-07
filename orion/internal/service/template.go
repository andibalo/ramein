package service

import "github.com/andibalo/ramein/orion/internal/repository"

type templateService struct {
	templateRepo repository.TemplateRepository
}

func NewTemplateService(templateRepo repository.TemplateRepository) *templateService {

	return &templateService{
		templateRepo: templateRepo,
	}
}

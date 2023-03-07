package service

import (
	"github.com/andibalo/ramein/orion/ent"
	"github.com/andibalo/ramein/orion/internal/config"
	"github.com/andibalo/ramein/orion/internal/repository"
	"github.com/andibalo/ramein/orion/internal/request"
	"go.uber.org/zap"
)

type templateService struct {
	cfg          config.Config
	templateRepo repository.TemplateRepository
}

func NewTemplateService(cfg config.Config, templateRepo repository.TemplateRepository) *templateService {

	return &templateService{
		cfg:          cfg,
		templateRepo: templateRepo,
	}
}

func (s *templateService) CreateTemplate(data request.CreateTemplateReq) error {

	err := s.templateRepo.Save(data)

	if err != nil {
		s.cfg.Logger().Error("[CreateTemplate] Error inserting template to db", zap.Error(err))
		return err
	}

	return nil
}

func (s *templateService) GetTemplateByID(templateID string) (*ent.Template, error) {

	template, err := s.templateRepo.GetByID(templateID)
	if err != nil {
		s.cfg.Logger().Error("[CreateTemplate] Error inserting template to db", zap.Error(err))
		return nil, err
	}

	return template, nil
}

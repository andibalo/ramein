package repository

import (
	"context"
	"github.com/andibalo/ramein/orion/ent"
	"github.com/andibalo/ramein/orion/ent/template"
	"github.com/andibalo/ramein/orion/internal/apperr"
	"github.com/andibalo/ramein/orion/internal/request"
	"github.com/google/uuid"
)

type templateRepository struct {
	db *ent.Client
}

func NewTemplateRepository(db *ent.Client) *templateRepository {

	return &templateRepository{
		db: db,
	}
}

func (r *templateRepository) Save(data request.CreateTemplateReq) error {

	_, err := r.db.Template.
		Create().
		SetName(data.Name).
		SetType(data.Type).
		SetTemplate(data.Template).
		SetCreatedBy("test"). //TODO: Remove later
		Save(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *templateRepository) GetByTemplateName(templateName string) (*ent.Template, error) {

	t, err := r.db.Template.
		Query().
		Where(template.Name(templateName)).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperr.ErrNotFound
		}

		return nil, err
	}

	return t, nil
}

func (r *templateRepository) GetByID(templateID string) (*ent.Template, error) {
	uid, _ := uuid.Parse(templateID)

	t, err := r.db.Template.
		Query().
		Where(template.ID(uid)).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperr.ErrNotFound
		}

		return nil, err
	}

	return t, nil
}

package repository

import (
	"context"
	"github.com/andibalo/ramein/orion/ent"
	"github.com/andibalo/ramein/orion/internal/request"
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

package repository

import "entgo.io/ent/entc/integration/ent"

type templateRepository struct {
	db *ent.Client
}

func NewTemplateRepository(db *ent.Client) *templateRepository {

	return &templateRepository{
		db: db,
	}
}

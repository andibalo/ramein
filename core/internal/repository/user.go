package repository

import (
	"context"
	"fmt"
	"github.com/andibalo/ramein/core/internal/model"
	"github.com/uptrace/bun"
)

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(user *model.User) error {

	res, err := r.db.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	fmt.Print(res)

	return nil
}

func (r *userRepository) GetByID(userID string) (*model.User, error) {
	user := &model.User{}

	err := r.db.NewSelect().Model(user).Where("id = ?", userID).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}

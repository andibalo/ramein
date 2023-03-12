package repository

import (
	"context"
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

func (r *userRepository) SaveTx(user *model.User, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Save(user *model.User) error {

	_, err := r.db.NewInsert().Model(user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}

	err := r.db.NewSelect().Model(user).Where("email = ?", email).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) SaveUserVerifyEmailTx(userVerifyEmail *model.UserVerifyEmail, tx bun.Tx) error {

	_, err := tx.NewInsert().Model(userVerifyEmail).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

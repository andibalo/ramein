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

func (r *userRepository) GetUserVerifyEmailByID(id string) (*model.UserVerifyEmail, error) {
	userVerifyEmail := &model.UserVerifyEmail{}

	err := r.db.NewSelect().Model(userVerifyEmail).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return userVerifyEmail, nil
}

func (r *userRepository) SetUserToEmailVerifiedTx(id string, tx bun.Tx) error {
	user := &model.User{}
	user.IsEmailVerified = true

	_, err := tx.NewUpdate().
		Model(user).
		Column("is_email_verified").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) SetUserVerifyEmailToUsed(id string, tx bun.Tx) error {
	userVerifyEmail := &model.UserVerifyEmail{}
	userVerifyEmail.IsUsed = true

	_, err := tx.NewUpdate().
		Model(userVerifyEmail).
		Column("is_used").
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

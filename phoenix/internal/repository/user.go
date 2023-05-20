package repository

import (
	"context"
	"github.com/andibalo/ramein/phoenix/internal/config"
	"github.com/andibalo/ramein/phoenix/internal/model"
	"github.com/andibalo/ramein/phoenix/internal/request"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type userRepo struct {
	ctx context.Context
	cfg config.Config
	db  neo4j.DriverWithContext
}

func NewUserRepo(ctx context.Context, cfg config.Config, db neo4j.DriverWithContext) *userRepo {
	return &userRepo{
		ctx: ctx,
		cfg: cfg,
		db:  db,
	}
}

func (r *userRepo) FetchUsers(req request.GetUsersListReq) ([]model.User, error) {

	session := r.db.NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: r.cfg.DbUserName()})
	defer session.Close(r.ctx)

	if req.Limit < 1 {
		req.Limit = 10
	}

	users := []model.User{}

	_, err := session.ExecuteRead(r.ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		query := `
			MATCH (user:User)
			RETURN user 
			LIMIT $limit`

		result, err := tx.Run(r.ctx, query, map[string]any{
			"limit": req.Limit,
		})
		if err != nil {
			return nil, err
		}

		for result.Next(r.ctx) {

			var user model.User

			u := result.Record().Values[0].(neo4j.Node)
			user.ID = u.GetId()
			user.FirstName = u.Props["firstName"].(string)
			user.LastName = u.Props["lastName"].(string)
			user.CoreUserID = u.Props["coreUserId"].(string)
			user.Gender = u.Props["gender"].(string)
			user.Email = u.Props["email"].(string)
			user.PhoneNumber = u.Props["phoneNumber"].(string)

			users = append(users, user)
		}
		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}

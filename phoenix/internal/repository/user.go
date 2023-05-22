package repository

import (
	"context"
	"fmt"
	"github.com/andibalo/ramein/phoenix/internal/apperr"
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

//func (r *userRepo) SaveUser(user model.User) error {
//
//	session := r.db.NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: r.cfg.DbUserName()})
//	defer session.Close(r.ctx)
//
//	_, err := session.ExecuteWrite(r.ctx, func(tx neo4j.ManagedTransaction) (any, error) {
//
//		query := `
//			MERGE (user:USER{})
//			RETURN user`
//
//		result, err := tx.Run(r.ctx, query, map[string]any{
//			"requested_at": neo4j.OffsetTimeOf(time.Now()),
//		})
//
//		if err != nil {
//			return nil, err
//		}
//
//		return nil, result.Err()
//	})
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

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

func (r *userRepo) FetchUserByUserID(coreUserID string) (model.User, error) {

	session := r.db.NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: r.cfg.DbUserName()})
	defer session.Close(r.ctx)

	user := model.User{}

	_, err := session.ExecuteRead(r.ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		query := `
			MATCH (user:User)
			WHERE user.coreUserId = $core_user_id
			RETURN user
			LIMIT 1`

		result, err := tx.Run(r.ctx, query, map[string]any{
			"core_user_id": coreUserID,
		})
		if err != nil {
			return nil, err
		}

		if !result.Next(r.ctx) {
			return nil, apperr.ErrNotFound
		}

		for result.Next(r.ctx) {

			u := result.Record().Values[0].(neo4j.Node)
			fmt.Println(u, "U")
			user.ID = u.GetId()
			user.FirstName = u.Props["firstName"].(string)
			user.LastName = u.Props["lastName"].(string)
			user.CoreUserID = u.Props["coreUserId"].(string)
			user.Gender = u.Props["gender"].(string)
			user.Email = u.Props["email"].(string)
			user.PhoneNumber = u.Props["phoneNumber"].(string)

		}
		return nil, result.Err()
	})

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepo) FetchFriendsListByUserID(userID string, req request.GetFriendsListReq) ([]model.User, error) {

	session := r.db.NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: r.cfg.DbUserName()})
	defer session.Close(r.ctx)

	if req.Limit < 1 {
		req.Limit = 10
	}

	users := []model.User{}

	_, err := session.ExecuteRead(r.ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		query := `
			MATCH (userFriends:User)-[:IS_FRIENDS_WITH]->(user:User{coreUserId: $core_user_id})
			RETURN userFriends 
			LIMIT $limit`

		result, err := tx.Run(r.ctx, query, map[string]any{
			"core_user_id": userID,
			"limit":        req.Limit,
		})
		if err != nil {
			return nil, err
		}

		for result.Next(r.ctx) {

			var userFriend model.User

			u := result.Record().Values[0].(neo4j.Node)
			userFriend.ID = u.GetId()
			userFriend.FirstName = u.Props["firstName"].(string)
			userFriend.LastName = u.Props["lastName"].(string)
			userFriend.CoreUserID = u.Props["coreUserId"].(string)
			userFriend.Gender = u.Props["gender"].(string)
			userFriend.Email = u.Props["email"].(string)
			userFriend.PhoneNumber = u.Props["phoneNumber"].(string)

			users = append(users, userFriend)
		}
		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) SaveFriendRequestRelationship(userID, targetUserID string) error {

	session := r.db.NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: r.cfg.DbUserName()})
	defer session.Close(r.ctx)

	_, err := session.ExecuteWrite(r.ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		query := `
			MATCH (user:User{coreUserId : $core_user_id})
			MATCH (targetUser:User{coreUserId : $target_core_user_id})
			MERGE (user)-[:HAS_FRIEND_REQUESTED{requestedAt: datetime(), hasBeenAccepted: false}]->(targetUser)
			RETURN user,targetUser`

		result, err := tx.Run(r.ctx, query, map[string]any{
			"core_user_id":        userID,
			"target_core_user_id": targetUserID,
		})

		if err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) SaveIsFriendsWithRelationship(userID, targetUserID string) error {

	session := r.db.NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: r.cfg.DbUserName()})
	defer session.Close(r.ctx)

	_, err := session.ExecuteWrite(r.ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		query := `
			MATCH (user:User{coreUserId : $core_user_id})-[friendReq:HAS_FRIEND_REQUESTED{hasBeenAccepted : false}]->(targetUser:User{coreUserId : $target_core_user_id})
			CREATE (user)-[:IS_FRIENDS_WITH{createdAt: datetime()}]->(targetUser)
			CREATE (targetUser)-[:IS_FRIENDS_WITH{createdAt: datetime()}]->(user)
			SET friendReq.hasBeenAccepted = true
			RETURN user,targetUser`

		result, err := tx.Run(r.ctx, query, map[string]any{
			"core_user_id":        userID,
			"target_core_user_id": targetUserID,
		})

		if err != nil {
			return nil, err
		}

		if !result.Next(r.ctx) {
			return nil, apperr.ErrNotFound
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) CheckIsFriendsWithRelationshipExist(userID, targetUserID string) (bool, error) {

	session := r.db.NewSession(r.ctx, neo4j.SessionConfig{DatabaseName: r.cfg.DbUserName()})
	defer session.Close(r.ctx)

	isFriends := false

	_, err := session.ExecuteWrite(r.ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		query := `
			MATCH (user:User{coreUserId : $core_user_id})-[r:IS_FRIENDS_WITH]->(targetUser:User{coreUserId : $target_core_user_id})
			RETURN user,targetUser`

		result, err := tx.Run(r.ctx, query, map[string]any{
			"core_user_id":        userID,
			"target_core_user_id": targetUserID,
		})

		if err != nil {
			return isFriends, err
		}

		if !result.Next(r.ctx) {
			return isFriends, apperr.ErrNotFound
		}

		isFriends = true

		return isFriends, result.Err()
	})

	if err != nil {
		return isFriends, err
	}

	return isFriends, nil
}

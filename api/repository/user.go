package repository

import (
	"context"
	"oauthive/db/entities"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
)

type UserRepository struct {
	db rel.Repository
}

func NewUserRepository(db rel.Repository) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) UpsertUser(ctx context.Context, user *entities.User) error {
	var existingUser entities.User
	err := r.db.Find(ctx, &existingUser, where.Eq("email", user.Email))
	if err != nil {
		if err == rel.ErrNotFound {
			return r.db.Insert(ctx, user)
		}
		return err
	}
	existingUser.Name = user.Name
	return r.db.Update(ctx, &existingUser)
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User
	err := r.db.Find(ctx, &user, where.Eq("email", email))
	if err != nil {
		if err == rel.ErrNotFound {
			return user, nil
		}
		return user, err
	}
	return user, nil
}

func (r *UserRepository) FindUserBySessionID(ctx context.Context, sessionID int) (entities.User, error) {
	var (
		session entities.Session
		user    entities.User
	)

	err := r.db.Find(ctx, &session, where.Eq("id", sessionID))
	if err != nil {
		return user, err
	}

	err = r.db.Find(ctx, &user, where.Eq("id", session.UserID))
	if err != nil {
		return user, err
	}

	return user, nil
}

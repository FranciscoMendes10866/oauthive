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

func (self *UserRepository) UpsertUser(ctx context.Context, user *entities.User) error {
	var existingUser entities.User
	err := self.db.Find(ctx, &existingUser, where.Eq("email", user.Email))
	if err != nil {
		if err == rel.ErrNotFound {
			return self.db.Insert(ctx, user)
		}
		return err
	}
	existingUser.Name = user.Name
	existingUser.Image = user.Image
	return self.db.Update(ctx, &existingUser)
}

func (self *UserRepository) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := self.db.Find(ctx, &user, where.Eq("email", email))
	if err != nil {
		if err == rel.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

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

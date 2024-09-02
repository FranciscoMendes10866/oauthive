package repository

import (
	"context"
	"oauthive/db/entities"

	"github.com/go-rel/rel"
)

type AccountRepository struct {
	db rel.Repository
}

func NewAccountRepository(db rel.Repository) *AccountRepository {
	return &AccountRepository{db}
}

func (self *AccountRepository) UpsertAccount(ctx context.Context, account *entities.Account) error {
	existingAccount := &entities.Account{}
	err := self.db.Find(
		ctx,
		existingAccount,
		rel.Eq("provider", account.Provider).
			AndEq("user_id", account.UserID),
	)

	if err != nil {
		if err == rel.ErrNotFound {
			err = self.db.Insert(ctx, account)
			if err != nil {
				return err
			}
		}
		return err
	}

	account.ID = existingAccount.ID
	return self.db.Update(ctx, account)
}

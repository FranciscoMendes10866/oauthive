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

func (repo *AccountRepository) UpsertAccount(ctx context.Context, account *entities.Account) error {
	existingAccount := &entities.Account{}
	err := repo.db.Find(ctx, existingAccount, rel.Eq("provider", account.Provider).AndEq("user_id", account.UserID))
	if err == rel.ErrNotFound {
		return repo.db.Insert(ctx, account)
	}
	if err != nil {
		return err
	}
	account.ID = existingAccount.ID
	return repo.db.Update(ctx, account)
}

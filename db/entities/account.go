package entities

import (
	"time"
)

type Account struct {
	ID                int       `json:"id" rel:"type:text;primary_key"`
	UserID            int       `json:"user_id" rel:"type:text;foreign_key:user;on_delete:cascade"`
	Provider          string    `json:"provider" rel:"type:text"`
	ProviderAccountID string    `json:"provider_account_id" rel:"type:text"`
	CreatedAt         time.Time `json:"created_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"updated_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (self *Account) Table() string {
	return "account"
}

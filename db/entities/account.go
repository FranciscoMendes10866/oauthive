package entities

import (
	"time"
)

type Account struct {
	ID                int       `json:"id" rel:"type:integer;primary_key"`
	UserID            int       `json:"user_id" rel:"type:integer;foreign_key:users;on_delete:cascade"`
	Provider          string    `json:"provider" rel:"type:text"`
	ProviderAccountID string    `json:"provider_account_id" rel:"type:text"`
	CreatedAt         time.Time `json:"created_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `json:"updated_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

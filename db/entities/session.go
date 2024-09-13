package entities

import (
	"time"
)

type Session struct {
	ID        int       `json:"id" rel:"type:integer;primary_key"`
	UserID    int       `json:"user_id" rel:"type:integer;foreign_key:user;on_delete:cascade"`
	ExpiresAt int64     `json:"expires_at" rel:"type:integer"`
	CreatedAt time.Time `json:"created_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

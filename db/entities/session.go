package entities

import (
	"time"
)

type Session struct {
	ID        int       `json:"id" rel:"type:text;primary_key"`
	UserID    int       `json:"user_id" rel:"type:text;foreign_key:user;on_delete:cascade"`
	Expires   time.Time `json:"expires" rel:"type:timestamp"`
	CreatedAt time.Time `json:"created_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (self *Session) Table() string {
	return "session"
}

package entities

import (
	"time"
)

type User struct {
	ID        int       `json:"id" rel:"type:integer;primary_key"`
	Name      string    `json:"name" rel:"type:text"`
	Email     string    `json:"email" rel:"type:text;unique"`
	CreatedAt time.Time `json:"created_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

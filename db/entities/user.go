package entities

import (
	"time"
)

type User struct {
	ID        int       `json:"id" rel:"type:text;primary_key"`
	Name      string    `json:"name" rel:"type:text"`
	Email     string    `json:"email" rel:"type:text;unique"`
	Image     string    `json:"image" rel:"type:text"`
	CreatedAt time.Time `json:"created_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" rel:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (self *User) Table() string {
	return "user"
}
